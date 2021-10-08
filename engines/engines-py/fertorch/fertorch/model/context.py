# COPIED FROM mxnet-model-server

class Context(object):
    """
    Context stores model relevant worker information
    Some fixed during load times and some
    """

    def __init__(self, model_name, model_dir, manifest, batch_size, gpu, mms_version):
        self.model_name = model_name
        self.manifest = manifest
        self._system_properties = {
            "model_dir": model_dir,
            "gpu_id": gpu,
            "batch_size": batch_size,
            "server_name": "MMS",
            "server_version": mms_version
        }

    @property
    def system_properties(self):
        return self._system_properties

    def __eq__(self, other):
        return isinstance(other, Context) and self.__dict__ == other.__dict__
