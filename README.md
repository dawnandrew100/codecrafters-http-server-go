[![progress-banner](https://backend.codecrafters.io/progress/http-server/020cc44e-dfb9-40d0-8daa-1b9c55fde346)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Message from Code Crafters

This is a starting point for Go solutions to the
["Build Your Own HTTP server" Challenge](https://app.codecrafters.io/courses/http-server/overview).

[HTTP](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol) is the
protocol that powers the web. In this challenge, you'll build a HTTP/1.1 server
that is capable of serving multiple clients.

Along the way you'll learn about TCP servers,
[HTTP request syntax](https://www.w3.org/Protocols/rfc2616/rfc2616-sec5.html),
and more.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# My goals with the project 

Go is a relatively new langauge for me, so I intend to complete this challenge (and the extension) to bolster my skills in the language!
An HTTP server is a good chance to learn not only basic syntax but unique ways of parsing and best practices when it comes to programming in go!

The hope is that by the end of the free challenge window (May 31st, 2024), I'll have a fully fledged Go server with compression!

As of today (May 29th, 2024), the server has been created through the base challenge and the gzip header can be read!

## Update!
Today (May 29th, 2024 @ 21:50) marks the end of my challenge! 
Full server features:
- Binds to port (currently port 4221)
- Responds with both success (200 and 201) and error (404 and 400) codes
- Extracts target path of request to perform some action
- Echos text when echo command is given
- Responds to both GET and POST requests
- Reads header (currently used for User-Agent)
- Handles concurrent connections
- Ignores erroneous compression requests (requests other than gzip)
- Compresses files using gzip
