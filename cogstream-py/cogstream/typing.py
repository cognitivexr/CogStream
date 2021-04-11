# from: https://github.com/thrau/pymq/blob/master/pymq/typing.py
# MIT License
#
# Copyright (c) 2019 Thomas Rausch
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
import inspect
import types
import typing
from pydoc import locate

import numpy as np

if hasattr(typing, '_GenericAlias'):
    def _is_generic(cls):
        return isinstance(cls, typing._GenericAlias)

else:
    if hasattr(typing, '_Union'):
        def _is_generic(cls):
            return isinstance(cls, typing.GenericMeta)

    else:
        def _is_generic(cls):
            raise RuntimeError('Need python>=3.6')


def is_generic(cls):
    """
    Detects any kind of generic, for example `List` or `List[int]`. This includes "special" types like
    Union and Tuple - anything that's subscriptable, basically.
    """
    return _is_generic(cls)


def load_class(classname):
    return locate(classname)


def new_instance(cls, data):
    # if available, use constructor args
    arg_names = inspect.getfullargspec(cls).args
    args = {k: v for k, v in data.items() if k in arg_names}
    if args:
        obj = cls(**args)
    else:
        obj = cls()

    # set all others via 'setattr'
    for key, value in data.items():
        if key in args:
            continue  # already set through constructor
        setattr(obj, key, value)

    return obj


def fullname(o):
    # o.__module__ + "." + o.__class__.__qualname__ is an example in
    # this context of H.L. Mencken's "neat, plausible, and wrong."
    # Python makes no guarantees as to whether the __module__ special
    # attribute is defined, so we take a more circumspect approach.
    # Alas, the module name is explicitly excluded from __qualname__
    # in Python 3.

    if isinstance(o, (types.MethodType, types.FunctionType)):
        return o.__module__ + '.' + o.__qualname__

    if isinstance(o, type):
        o = o
    else:
        o = o.__class__

    module = o.__module__
    if module is None or module == str.__class__.__module__:
        return o.__name__  # Avoid reporting __builtin__
    else:
        return module + '.' + o.__name__


def deep_from_dict(doc, cls):
    if doc is None:
        return doc

    if type(doc) == cls:
        return doc

    if cls == typing.Any:
        return doc

    if cls == type:
        raise TypeError('Deserializing types is not safe')

    if is_generic(cls):
        container_class = cls.__origin__

        if issubclass(container_class, list):
            element_class = cls.__args__[0]
            return [deep_from_dict(element, element_class) for element in doc]

        if issubclass(container_class, set):
            element_class = cls.__args__[0]
            return {deep_from_dict(element, element_class) for element in doc}

        if issubclass(container_class, tuple):
            return tuple([deep_from_dict(doc[i], cls.__args__[i]) for i in range(len(doc))])

        if issubclass(container_class, dict):
            key_type = cls.__args__[0]
            value_type = cls.__args__[1]
            return {deep_from_dict(k, key_type): deep_from_dict(v, value_type) for k, v in doc.items()}

        raise TypeError('Unknown generic class %s' % cls)

    if issubclass(cls, Exception):
        if isinstance(doc, (list, tuple)):
            return cls(*doc)
        else:
            return cls(doc)

    if isinstance(doc, (bool, int, float, str, bytes, bytearray)):
        if type(doc) != cls:
            return cls(doc)
        return doc

    if isinstance(doc, list) and cls in (set, tuple):
        return cls(doc)

    # otherwise we treat it as an object
    spec = typing.get_type_hints(cls)
    result = dict()

    if isinstance(doc, (list, tuple)):
        # named tuples for example may be in a list
        i = 0
        for name, target_type in spec.items():
            result[name] = deep_from_dict(doc[i], target_type)
            i += 1

    else:
        for name, target_type in spec.items():
            if name not in doc:
                continue
            result[name] = deep_from_dict(doc[name], target_type)

    return new_instance(cls, result)


def deep_to_dict(obj):
    if obj is None:
        return None

    if isinstance(obj, (bool, int, float, str, bytes, bytearray)):
        return obj

    if isinstance(obj, np.number):
        return obj.item()

    if isinstance(obj, tuple):
        return tuple([deep_to_dict(a) for a in obj])

    if isinstance(obj, list):
        return [deep_to_dict(a) for a in obj]

    if isinstance(obj, dict):
        return {k: deep_to_dict(v) for k, v in obj.items()}

    if isinstance(obj, set):
        return {deep_to_dict(a) for a in obj}

    if isinstance(obj, (type, types.MethodType, types.FunctionType)):
        return fullname(obj)

    if isinstance(obj, Exception):
        return deep_to_dict(obj.args)

    if hasattr(obj, '__dict__'):
        return deep_to_dict(obj.__dict__)

    raise TypeError('Unhandled type %s' % type(obj))
