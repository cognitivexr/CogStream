import os
import glob

import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

with open("requirements-dev.txt", "r") as fh:
    tests_require = [line for line in fh.read().split(os.linesep) if line]

with open("requirements.txt", "r") as fh:
    install_requires = [line for line in fh.read().split(os.linesep) if line]

setuptools.setup(
    name="cogstream-py",
    version="0.1.0",
    author="Thomas Rausch",
    author_email="thomas@cognitivexr.at",
    description="CogStream: AI-based video analytics system",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/cognitivexr/cogstream",
    packages=setuptools.find_packages(),
    setup_requires=['wheel'],
    test_suite="tests",
    tests_require=tests_require,
    install_requires=install_requires,
    pyton_requires='>=3.6',
    classifiers=[
        "Programming Language :: Python :: 3",
        "Operating System :: OS Independent",
    ],
    entry_points={
    },

)
