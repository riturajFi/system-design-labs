# system-design-labs

Hands-on implementations of real-world system design problems in **Go** and **Rust**.
This repository focuses on building **production-oriented backend and distributed systems** from first principles, with emphasis on correctness, concurrency, failure handling, and observability.

---

## Goals

* Convert system design concepts into **working code**
* Understand **trade-offs**, not just final architectures
* Practice **incremental design → implementation**
* Build intuition for **scalability, reliability, and performance**

This is **not** an interview cheat sheet. Every project is executable, testable, and stress-tested where applicable.

---

## What You’ll Find Here

Each system is implemented as an independent module and typically includes:

* Problem statement & requirements
* Core design decisions and assumptions
* Go and/or Rust implementation
* Concurrency and synchronization logic
* Basic tests and load/stress scenarios
* Notes on limitations and next improvements

Example systems:

* ID generators
* URL shorteners
* Rate limiters
* Job queues
* Web crawlers
* Caching layers

---

## Repository Structure (WIP)

```text
/
├── go/
│   ├── id-generator/
│   ├── url-shortener/
│   └── ...
├── rust/
│   ├── id-generator/
│   ├── rate-limiter/
│   └── ...
└── docs/
    └── design-notes/
```

---

## Philosophy

* Prefer **simple, explicit code** over clever abstractions
* Optimize for **clarity before performance**
* Measure behavior under load instead of guessing
* Treat failures as first-class design inputs

---

## Status

🚧 Actively evolving.
Systems will be added incrementally and refined over time.
