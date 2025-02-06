# Imagika - High-Performance Image Proxy Server

![Imagika Demo](cmd/assets/imagika-preview.gif)

Imagika is a high-performance image proxy server designed to optimize the delivery of high-resolution images (e.g., 4K images) by leveraging caching, rate limiting, and efficient request handling. The project consists of a frontend, an origin server where high-resolution images are stored, and a proxy server that handles image requests, caching, and delivery.

---

## Working of Proxy

### Overview

Imagika is built with a client-server architecture:

1. **Frontend**: The user-facing interface where users can request images.
2. **Origin Server**: Stores the original high-resolution (4K) images.
3. **Proxy Server**: Acts as an intermediary between the frontend and the origin server. It handles image requests, applies caching, and ensures efficient delivery.

### Key Components

- **Frontend**: Built with a modern web framework (e.g., React, Vue, or similar), the frontend allows users to request images by sending requests to the proxy server.
- **Origin Server**: A server or storage system (e.g., AWS S3, local storage) that hosts the original high-resolution images.
- **Proxy Server**: The core of Imagika, located at `/cmd/proxy`, handles incoming requests, fetches images from the origin server, and caches them for faster delivery.

### Proxy Logic

The proxy server is responsible for:

1. **Request Handling**: Receives image requests from the frontend.
2. **Caching**: Uses an **LRU (Least Recently Used) cache** to store frequently requested images in memory, reducing the need to fetch them repeatedly from the origin server.
3. **Rate Limiting**: Implements rate limiting to prevent abuse and ensure fair usage of the proxy server.
4. **Thread Safety**: The LRU cache is thread-safe, ensuring that multiple concurrent requests do not lead to race conditions or data corruption.
5. **Image Delivery**: Delivers the requested image to the frontend, either from the cache or by fetching it from the origin server.

---

## Features

- **LRU Caching**: Efficiently caches frequently requested images using a thread-safe LRU cache.
- **Rate Limiting**: Protects the server from abuse by limiting the number of requests from a single client.
- **High Performance**: Optimized for delivering high-resolution images with minimal latency.
- **Scalable**: Designed to handle a large number of concurrent requests.
- **Easy Setup**: Simple setup process for local development and testing.

---