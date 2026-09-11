# User Flows

Flows describe what the user does and what the system presents in response from the user's perspective.

## 1. First Run & Connecting to Linden

1. **Accessing the Interface:** The user opens any modern web browser on a machine connected to the same local network and navigates to `http://linden.local:8080` (or `http://localhost:8080` on the host machine).
2. **Initial Presentation:** The embedded SvelteKit web client loads immediately, presenting the chat workspace connected to the local inference backend with the default local model (`tinyllama`) active.

## 2. Sending a Message

1. **Input:** The user types a message in the chat input area at the bottom of the screen and presses `Enter` (or clicks the Send button).
2. **Immediate Feedback:** The user's message is immediately added to the conversation view, and the input field clears.
3. **Streaming Generation:** Linden sends the message to `/chat` and streams tokens back incrementally via Server-Sent Events. The user sees the assistant's reply render in real-time with full Markdown and syntax formatting.
4. **Completion:** When the inference backend completes the generation, the stream closes and the interface transitions back to ready status.

## 3. Recovering From a Backend Failure

1. **Inference Backend Offline:** If Ollama or the local inference engine is stopped or unreachable when a message is sent:
   - The user receives an immediate, clear error notification stating that the inference backend is unreachable.
   - The user's input is retained so it is not lost.
   - The system logs record the failure with an associated `X-Request-ID` without recording the user's message text.
2. **Empty Model Catalog or Missing Default Model:** If no models are pulled or installed in the backend:
   - When a user sends a prompt, the system returns a clear, actionable error indicating that the backend model is not installed.
   - The user is prompted to run `ollama pull tinyllama` (or the desired model) on the host machine.
