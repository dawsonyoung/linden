# The Linden Trust Manifesto

Linden was built on a core philosophy: **Honesty and Transparency first.**

In an era where personal assistants operate primarily in the cloud, processing your most sensitive thoughts and files on corporate servers, Linden takes a different approach. We are building a local-first appliance that puts you in control.

This document outlines our current stance, our limitations, and our long-term promises regarding your privacy and safety.

## 1. Transparency and Our Current Limitations

Linden is designed as a LAN-bound, local-first appliance. However in this early stage of development, the current architecture relies on third-party dependencies (like inference engines) that we do not fully control. Until we have absolute control over these underlying dependencies and have implemented our own foundational, privacy-focused features, we cannot make absolute guarantees. 

Our promise to you today is **transparency**. We promise to always be honest when we make changes to our privacy structures and to make transparency paramount as we develop this tool.

For example, Linden allows you to configure custom endpoint URLs for your inference engine (such as Ollama). If you manually point Linden to a cloud-hosted instance or an external API provider, **Linden cannot control what happens to your data once it leaves your local network.** Routing your chats to a third-party cloud provider transfers the privacy responsibility to that provider.

## 2. Our Privacy Commitment for Modern Capabilities

We put privacy first. That means when we introduce features that connect to the outside world, we adhere to these strict rules:
* **Auditable:** Any data that would potentially leave the Linden system will be fully auditable by you, the owner.
* **Optional:** External connections will always be strictly optional. You can reject them entirely and still use the core offline features.
* **Plain English:** We will explain exactly what data is moving, where it is going, and why, in plain English—no confusing legal jargon.

## 3. Safety and Data Provenance

As we build Linden, safety is just as important as privacy. 

We intend to develop features that will give you total control over Linden's behavior. A key part of this safety is **data provenance**—ensuring that the answers Linden gives you are grounded in reality.

We are building towards a future where Linden's responses are based strictly on **real, local data** that you have provided (your documents, your notes, your facts), rather than relying on the assumptions, hallucinations, or amalgamations gathered from the internet at large by the underlying AI models. You should always know exactly *why* Linden gave you an answer and *where* that answer came from.

***

We build tools, not authorities. As Linden evolves, our promise is to continuously transfer total control into your hands, backed by uncompromising honesty and transparency.
