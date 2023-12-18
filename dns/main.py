import socket


def build_dns_query(domain: str):
    header = bytearray([
        0x00, 0x01,  # Transaction ID
        0x00, 0x00,  # Flags: Standard query
        0x00, 0x01,  # Questions
        0x00, 0x00,  # Answer RRs
        0x00, 0x00,  # Authority RRs
        0x00, 0x00   # Additional RRs
    ])

    question = bytearray()
    labels = domain.split('.')
    for label in labels:
        question.append(len(label))
        question.extend(label.encode('utf-8'))
    question.extend([0x00, 0x00, 0x01, 0x00, 0x01])  # QTYPE and QCLASS (A record, Internet)

    return header + question


def send_dns_query(query, server, port=2053):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
        s.sendto(query, (server, port))
        response, _ = s.recvfrom(1024)
    return response


def parse_dns_response(response):
    print(response)
    print(response.hex())


if __name__ == "__main__":
    dns_server = "127.0.0.1"
    domain = "habr.com"
    dns_query = build_dns_query(domain)
    dns_response = send_dns_query(dns_query, dns_server)
    parse_dns_response(dns_response)
