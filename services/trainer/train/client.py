import os
import random
import socket


def main():
    host_id = get_host_id()
    print(host_id)


def parse_hostname():
    hostname = socket.gethostname()
    doc_domain = ".doc.ic.ac.uk"
    if doc_domain not in hostname:
        return None
    return hostname.replace(doc_domain, '')


def get_host_id():
    hostname = parse_hostname()
    if hostname:
        return hostname
    return generate_host_id()


def generate_host_id():
    integer_id = str(random.randint(1, 10000))
    host_id = "vast" + integer_id
    os.environ['HOST_ID'] = host_id
    return host_id


if __name__ == "__main__":
    main()
