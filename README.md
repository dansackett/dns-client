# DNS Client

This repository hosts code for a simple DNS implementation. This project is not
meant for production use but rather it is how I learned how DNS works under the
hood. The code itself offers a functional DNS client (like `dig`) with a very
limited feature set (You can only query for record types from a DNS server).

I'm hoping to continue learning about different resource records and then move
into learning about DNS servers themselves. The message encoding and decoding
can be used for a server application for example.

The following information about DNS is what I have learned so far.

## Using this

If you want to try this project out, you can either do `go get
github.com/dansackett/dns-client` or you can download the source here and run
`go build`.

You can then use the application with the following arguments:

```
$ ./dns-client -help
Usage of ./dns:
  -domain string
        The domain to run DNS queries on. This is required.
  -server-addr string
        IP and Port for the DNS server to query. Defaults to "8.8.8.8:53". (default "8.8.8.8:53")
  -type string
        Record type to lookup. Defaults to "A" (default "A")
```

## What is DNS?

DNS (Domain Name System)(Domain Name System) is one of the core features of the
internet. Without it, we wouldn't know that github.com routes to the IP Address
`192.30.253.112`. You see, behind the fancy domain names are 4 bytes which make
up an IP address.

You can think of DNS as a phonebook. A phone number is something recognizable
and when you dial it you get directed to a person. The domain name is our phone
number and it points us to an IP Address and subsequently a website.

Every internet device has a unique IP address allowing other machines to find
them.

The force behind DNS are the DNS servers which interpret a domain name and do
alll of the magic to send you where you need to go.

## How does it work?

Very carefully...

[Cloudflare](https://www.cloudflare.com/learning/dns/what-is-dns/) has a good
explanation which I'll paste here:

1. A user types 'example.com' into a web browser and the query travels into the
   Internet and is received by a DNS recursive resolver.
2. The resolver then queries a DNS root nameserver (.).  3 The root server then
   responds to the resolver with the address of a Top Level Domain (TLD) DNS
   server (such as .com or .net), which stores the information for its domains.
   When searching for example.com, our request is pointed toward the .com TLD.
4. The resolver then makes a request to the .com TLD.
5. The TLD server then responds with the IP address of the domain’s nameserver,
   example.com.
6. Lastly, the recursive resolver sends a query to the domain’s nameserver.
7. The IP address for example.com is then returned to the resolver from the
   nameserver.
8. The DNS resolver then responds to the web browser with the IP address of the
   domain requested initially.
9. The browser makes a HTTP request to the IP address.
10. The server at that IP returns the webpage to be rendered in the browser

There are a number of terms in there that might not make sense without a proper
introduction.

### Recursive Resolvers

These are servers which are designed to receive queries through clients like a
web browser. The recursor makes additional requests to satisfy the DNS query
being asked.

In short, it is the workhorse server and is the entrypoint to most DNS queries.
It will make multiple requests to different servers until it has the
information it needs to resolve the query.

Because this can be an intensive task to do, recursive resolvers generally use
caching systems in order to reduce queries and serve results quicker.

### Root Nameserver

The root nameserver is the first place a query is sent from the recursive
resolver. They handle a root zone file which either resolves a domain or points
to the Top Level Domain (TLD) server which would have information about the
domain being queried.

For example, say we're searching for `example.com` still. It likely won't know
the IP address of that however it will know the IP address of the TLD `.com`
server. It is on the `.com` server that a new request can be made for
`example.com`.

### Top Level Domain (TLD) Servers

These servers each host information about domains at the specific TLD being
queried. There are TLD servers for `.com, .net, .org, etc` and each one has
information about the domains registered with those TLDs.

### Authoratative nameserver

This is the last place a recursive resolved will look for an IP address. It is
the server which is actually holding the resource record for the domain you are
querying. It does not need to query any other sources since it is the source
that we have been looking for all along.

## What types of queries are made?

DNS queries generally come in three different types:

1. Recursive queries: The client making the request expects an answer in this
   case and doesn't mind how long it takes. In an ideal world, the information
   will be cached but this type of query would expect the recursive resolver
   server to make all the queries it needed to in order to find an answer.
2. Iterative queries: The client making the request is looking for the **best**
   answer they can get without recursion. In the case that a record cannot be
   found, the server returns a referral address to where we can look for the
   next step of the query process.
3. Non-recursive queries: This is usually only used when the client knows a
   domain is cached or it is querying an authoratative nameserver. It expects
   an answer because it is sure that the answer exists here.

## How does caching work?

Caching is the hero of most performance on the internet. Like in any
application, a cache is a collection of information in memory that can be
easily accessed for faster results. Items are added to a cache when they have
been queried and then any subsequent queries in a specific time frame making
the same request will be able to return the same value without looking it up
again.

DNS works this way in that if it recognizes a domain name that was queried
recently it can return the IP address immediately rather than recursively
looking for the information again.

DNS data is cached based on its TTL (time to live) and can be found in a couple
locations:

### Browser Cache

Web browsers have the ability to cache DNS records for you so you don't need to
make any requests outside of the browser session you're in. It can simply get
the request, check its cache, and return the IP to you immediately making the
pages served much faster.

### OS Cache

If the browser hasn't cached the results, a query will be sent to the operating
system before it leaves the machine for a recursive resolver. The OS has what
is called a **Stub Resolver** which is a DNS client. It has its own cache and if
the query can not be found in its cache it forwards the request to a recursive
resolver in their ISP.

## DNS Messages

While the process for answering DNS queries is interesting, it is all run by a
system of messages that are sent back and forth between servers. These messages
are defined in [RFC 1035](https://tools.ietf.org/html/rfc1035) and are made up
of multiple pieces of data. These messages are encoded in network byte order
and each client and server must understand how to interpret them.

### Message Header

The message header is included in all messages and sets the relevant flags
needed for the request or response. A message header is made up of the
following structure:

```
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      ID                       |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                    QDCOUNT                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                    ANCOUNT                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                    NSCOUNT                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                    ARCOUNT                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

Each row corresponds to two bytes (two octets). The fields included are:

- **ID:** A 16 bit identifier assigned by the program that generates any kind of
  query.  This identifier is copied the corresponding reply and can be used by
  the requester to match up replies to outstanding queries.
- **QR:**  A one bit field that specifies whether this message is a query (0), or
  a response (1).
- **OPCODE:** A four bit field that specifies kind of query in this message.
  This value is set by the originator of a query and copied into the response.
  The values are:
  - 0: a standard query (QUERY)
  - 1: an inverse query (IQUERY)
  - 2: a server status request (STATUS)
  - 3-15: reserved for future use
- **AA:** Authoritative Answer - this bit is valid in responses, and specifies
  that the responding name server is an authority for the domain name in
  question section.  Note that the contents of the answer section may have multiple owner
  names because of aliases.  The AA bit corresponds to the name which matches
  the query name, or the first owner name in the answer section.
- **TC:** TrunCation - specifies that this message was truncated due to length
  greater than that permitted on the transmission channel.
- **RD:** Recursion Desired - this bit may be set in a query and is copied into
  the response.  If RD is set, it directs the name server to pursue the query
  recursively.  Recursive query support is optional.
- **RA:** Recursion Available - this be is set or cleared in a response, and
  denotes whether recursive query support is available in the name server.
- **Z:** Reserved for future use.  Must be zero in all queries and responses.
- **RCODE:** Response code - this 4 bit field is set as part of responses.  The
  values have the following interpretation:
  - 0: No error condition
  - 1: Format error - The name server was unable to interpret the query.
  - 2: Server failure - The name server was unable to process this query due to a
    problem with the name server.
  - 3: Name Error - Meaningful only for responses from an authoritative name
    server, this code signifies that the domain name referenced in the query does
    not exist.
  - 4: Not Implemented - The name server does not support the requested kind of query.
  - 5: Refused - The name server refuses to perform the specified operation for
    policy reasons.  For example, a name server may not wish to provide the
    information to the particular requester, or a name server may not wish to
    perform a particular operation (e.g., zone transfer) for particular data.
  - 6-15: Reserved for future use.
- **QDCOUNT:** an unsigned 16 bit integer specifying the number of entries in
  the question section.
- **ANCOUNT:** an unsigned 16 bit integer specifying the number of resource
  records in the answer section.
- **NSCOUNT:** an unsigned 16 bit integer specifying the number of name server
  resource records in the authority records section.
- **ARCOUNT:** an unsigned 16 bit integer specifying the number of resource
  records in the additional records section.

### Message Questions

Questions are the actual queries that are being sent. You can send multiple
questions in a message and each one allows the responding server to return
information to fulfill those questions. THe question part of the message looks
like:

```
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                                               |
    /                     QNAME                     /
    /                                               /
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                     QTYPE                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                     QCLASS                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

Like the header, each row corresponds to two bytes (two octets). The fields
included are:

- **QNAME:** a domain name represented as a sequence of labels, where each
  label consists of a length octet followed by that number of octets.  The
  domain name terminates with the zero length octet for the null label of the
  root.  Note that this field may be an odd number of octets; no padding is
  used.
- **QTYPE:** a two octet code which specifies the type of the query.  The
  values for this field include all codes valid for a TYPE field, together with
  some more general codes which can match more than one type of RR.
- **QCLASS:** a two octet code that specifies the class of the query.  For
  example, the QCLASS field is IN for the Internet.

### Message Answers (Resource Records)

The answers to queries are called RRs or resource records. There are many
different types of RRs but they all have the same message format:

```
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                                               |
    /                                               /
    /                      NAME                     /
    |                                               |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      TYPE                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                     CLASS                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      TTL                      |
    |                                               |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                   RDLENGTH                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
    /                     RDATA                     /
    /                                               /
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

The fields are defined as:

- **NAME:** a domain name to which this resource record pertains.
- **TYPE:** two octets containing one of the RR type codes.  This field
  specifies the meaning of the data in the RDATA field.
- **CLASS:** two octets which specify the class of the data in the RDATA field.
- **TTL:** a 32 bit unsigned integer that specifies the time interval (in
  seconds) that the resource record may be cached before it should be
  discarded.  Zero values are interpreted to mean that the RR can only be used
  for the transaction in progress, and should not be cached.
- **RDLENGTH:** an unsigned 16 bit integer that specifies the length in octets
  of the RDATA field.
- **RDATA:** a variable length string of octets that describes the resource.
  The format of this information varies according to the TYPE and CLASS of the
  resource record.  For example, the if the TYPE is A and the CLASS is IN, the
  RDATA field is a 4 octet ARPA Internet address.

### Domain Names and Storage Considerations

Domains names are really a series of labels when broken down. For example, the
domain `something.example.com` where the labels are:

- `something`: a subdomain
- `example`: the hostname
- `com`: the TLD
- ` `: while not seen, every domain can be broken down to the root which is empty

Domain names are a prime candidate to compress and therefore there is a scheme
to do so. Pointers are used to indicate a different position in the bytes
returned to the client. When reading a domain name field you can check for the
following structure to determine if it's a pointer:

```
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    | 1  1|                OFFSET                   |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

If the two least significant bits are 1's then there is a pointer. The offset
tells you where you can find the domain name in the current message. At this
offset, the normal label structure is found.

Normal labels take the form of:

```
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    | 0  0|     SIZE        |      LABEL            |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

You can see that labels are identified by a `0` byte to start and then the
first byte is used to denote the size of the following label. Based on this
size you would now know how many bytes to read from the message before checking
for the next pointer or label instance.

Domain labels in a message are terminated by a `0` byte so you can know when to
move to the next field in the question or rr. As an example label, consider the
following domain names (from the RFC) F.ISI.ARPA, FOO.F.ISI.ARPA, ARPA, and the
root:

```
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    20 |           1           |           F           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    22 |           3           |           I           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    24 |           S           |           I           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    26 |           4           |           A           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    28 |           R           |           P           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    30 |           A           |           0           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    40 |           3           |           F           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    42 |           O           |           O           |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    44 | 1  1|                20                       |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    64 | 1  1|                26                       |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    92 |           0           |                       |
       +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

## Record Types

There are many different DNS record types. You can read about the current and
obsolete ones on [WikiPedia](https://en.wikipedia.org/wiki/List_of_DNS_record_types).

Some of the most common record types are:

- **A:** The host address (IP address)
- **NS:** An authoratative nameserver domain
- **CNAME:** Canonical name or alias for a domain
- **SOA:** Start of authority
- **MX:** Mail exchange information
- **AAAA:** IPV6 address
- **TXT:** Text strings (other hosts look for these sometimes to ensure authority)
