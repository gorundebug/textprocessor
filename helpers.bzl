# helpers.bzl

def get_basename(path):
    return path.split('/')[::-1][0]

def get_dirname(path, removeCount):
    components = path.split('/')
    if len(components) >= 1 + removeCount:
        return '/'.join(components[:-removeCount])
    else:
        return path

def get_current_dir_name():
    return get_basename(get_package_name())

def get_parent_dir_name():
    return get_basename(get_dirname(get_package_name(), 1))

def get_package_name_without(removeCount):
    return get_dirname(get_package_name(), removeCount)

def get_parent_package_name():
    return get_package_name_without(1)

def get_package_name():
    return native.package_name()

def repo_name(importpath):
    path_segments = importpath.split("/")
    segments = reversed(path_segments[0].split(".")) + path_segments[1:]
    candidate_name = "_".join(segments).replace("-", "_")
    return "".join([c.lower() if c.isalnum() else "_" for c in candidate_name.elems()])

def print_debug(message):
    print(message)