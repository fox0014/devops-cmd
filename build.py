#!/usr/bin/env python

from __future__ import print_function
import os
import time
import subprocess
import json
import platform
import sys
import ssl
import argparse
import warnings

WINDOWS_SAVE_BINARY = '-o bin/myutils.exe'
LINUX_SAVE_BINARY = '-o bin/myutils'


def _implementation():
    """Return a dict with the Python implementation and version.

    Provide both the name and the version of the Python implementation
    currently running. For example, on CPython 2.7.5 it will return
    {'name': 'CPython', 'version': '2.7.5'}.

    This function works best on CPython and PyPy: in particular, it probably
    doesn't work for Jython or IronPython. Future investigation should be done
    to work out the correct shape of the code for those platforms.
    """
    implementation = platform.python_implementation()

    if implementation == 'CPython':
        implementation_version = platform.python_version()
    elif implementation == 'PyPy':
        implementation_version = '%s.%s.%s' % (sys.pypy_version_info.major,
                                               sys.pypy_version_info.minor,
                                               sys.pypy_version_info.micro)
        if sys.pypy_version_info.releaselevel != 'final':
            implementation_version = ''.join([
                implementation_version, sys.pypy_version_info.releaselevel
            ])
    elif implementation == 'Jython':
        implementation_version = platform.python_version()  # Complete Guess
    elif implementation == 'IronPython':
        implementation_version = platform.python_version()  # Complete Guess
    else:
        implementation_version = 'Unknown'

    return {'name': implementation, 'version': implementation_version}


def info():
    """Generate information for a bug report."""
    try:
        platform_info = {
            'system': platform.system(),
            'release': platform.release(),
        }
    except IOError:
        platform_info = {
            'system': 'Unknown',
            'release': 'Unknown',
        }

    implementation_info = _implementation()

    system_ssl = ssl.OPENSSL_VERSION_NUMBER
    system_ssl_info = {
        'version': '%x' % system_ssl if system_ssl is not None else ''
    }

    return {
        'platform': platform_info,
        'implementation': implementation_info,
        'system_ssl': system_ssl_info,
    }


def py_check_compatibility():
    result = info()
    major, minor, patch = result['implementation']['version'].split('.')[:3]
    major, minor, patch = int(major), int(minor), int(patch)
    if major == 2:
        assert (2, 7, 0) <= (major, minor, patch)
    if major == 3:
        assert (3, 6, 0) <= (major, minor, patch)

# filter python version.
def py_version_filter():
    try:
        name = _implementation()['name']
        version = _implementation()['version']
        py_check_compatibility()
    except (AssertionError, ValueError):
        warnings.warn("python ({}) version ({})) doesn't match a supported "
                      "version! py2 > 2.7, py3 > 3.6 ".format(name, version))
        sys.exit(1)

# set linux build environ.
def set_linxbuild_environ():
    os.environ['CGO_ENABLED']="0"
    os.environ['GOARCH']="amd64"
    os.environ['GOOS']="linux"

# run command.
def runCmd(cmd):
    p = subprocess.Popen(cmd, shell = True, stdout = subprocess.PIPE, stderr = subprocess.PIPE)
    stdout = p.communicate()[0].decode('utf-8').strip()
    stderr = p.communicate()[1].decode('utf-8').strip()
    if stderr:
        print(stderr)
    return stdout
 
# Get last tag.
def lastTag():
    return runCmd('git describe --abbrev=0 --tags')
 
# Get current branch name.
def branch():
    return runCmd('git rev-parse --abbrev-ref HEAD')
 
# Get last git commit id.
def lastCommitId():
    return runCmd('git log --pretty=format:"%h" -1')
 
# Assemble build command.
def buildCmd(os):
    buildFlag = []
 
    version = lastTag()
    if version != "":
        buildFlag.append("-X 'myutils/myutils/cmd._version_={}'".format(version))
 
    branchName = branch()        
    if branchName != "":
        buildFlag.append("-X 'myutils/myutils/cmd._branch_={}'".format(branchName))
 
    commitId = lastCommitId()
    if commitId != "":
        buildFlag.append("-X 'myutils/myutils/cmd._commitId_={}'".format(commitId))

    # current time
    buildFlag.append("-X 'myutils/myutils/cmd._buildTime_={}'".format(time.strftime("%Y-%m-%d %H:%M %z")))
    build_args = 'go build -ldflags "-w -s {}" {}'.format(" ".join(buildFlag), WINDOWS_SAVE_BINARY)

    if os == 'linux':
        set_linxbuild_environ()
        build_args = 'go build -ldflags "-w -s {}" {}'.format(" ".join(buildFlag), LINUX_SAVE_BINARY)

    print(build_args)

    return build_args
 
def binary_compress():
    command = "upx.exe -9 .\\bin\*"
    print(runCmd(command))

def get_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-os", "--os", help="windos or linux")

    # 解析所传入的参数
    args = parser.parse_args()
    return(parser.parse_args())


def main():
    args = get_args()
    os = args.os
    py_version_filter()
    if subprocess.call(buildCmd(os), shell = True) == 0:
        print("build finished.")
        time.sleep(1)
        binary_compress()
    else:
        print("build failed.")


if __name__ == '__main__':
    main()
