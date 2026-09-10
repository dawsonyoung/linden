<script lang="ts">
	let messages = $state<{ role: string; content: string }[]>([]);
	let input = $state("");
	let isLoading = $state(false);
	let error = $state<Error | null>(null);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!input.trim() || isLoading) return;

		const userMsg = input;
		input = "";
		messages.push({ role: "user", content: userMsg });
		isLoading = true;
		error = null;

		const assistantMsgIndex = messages.length;
		messages.push({ role: "assistant", content: "" });

		try {
			// We send all messages except the one we just added for the assistant
			const payloadMessages = messages.slice(0, -1);
			const res = await fetch("/chat", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ model: "tinyllama", messages: payloadMessages }),
			});

			if (!res.ok) {
				const text = await res.text();
				let errMsg = text;
				try {
					const json = JSON.parse(text);
					if (json && json.error) errMsg = json.error;
				} catch {}
				throw new Error(errMsg || `Request failed with status ${res.status}`);
			}

			const reader = res.body?.getReader();
			if (!reader) throw new Error("No response body");

			const decoder = new TextDecoder();
			let buffer = "";
			let currentEvent = "";

			while (true) {
				const { done, value } = await reader.read();
				if (done) break;

				buffer += decoder.decode(value, { stream: true });
				const lines = buffer.split("\n");
				buffer = lines.pop() || "";

				for (const line of lines) {
					if (line.startsWith("event: ")) {
						currentEvent = line.slice(7).trim();
					} else if (line.startsWith("data: ")) {
						try {
							const data = JSON.parse(line.slice(6));
							if (currentEvent === "error" || data.error) {
								throw new Error(data.error || "An error occurred during generation");
							} else if (data.text) {
								messages[assistantMsgIndex].content +=
									data.text;
							}
						} catch (e: any) {
							if (currentEvent === "error" || (e.message && !e.message.startsWith("JSON"))) {
								throw e;
							}
						}
					} else if (line === "") {
						currentEvent = "";
					}
				}
			}
		} catch (err: any) {
			error = err;
		} finally {
			isLoading = false;
		}
	}
</script>

<div
	class="flex flex-col h-screen bg-base-100 text-base-content max-w-4xl mx-auto"
>
	<header class="p-4 border-b border-base-300">
		<h1 class="text-xl font-bold">Linden AI</h1>
	</header>

	<main class="flex-1 overflow-y-auto p-4 space-y-4 pb-32">
		{#each messages as msg}
			<div class="chat {msg.role === 'user' ? 'chat-end' : 'chat-start'}">
				<div class="chat-header capitalize mb-1 opacity-50">
					{msg.role}
				</div>
				<div
					class="chat-bubble {msg.role === 'user'
						? 'chat-bubble-primary'
						: 'chat-bubble-neutral'} whitespace-pre-wrap"
				>
					{msg.content}
				</div>
			</div>
		{/each}

		{#if isLoading && messages[messages.length - 1]?.role === "user"}
			<div class="chat chat-start">
				<div class="chat-bubble chat-bubble-neutral opacity-50">
					<span class="loading loading-dots loading-sm"></span>
				</div>
			</div>
		{/if}

		{#if error}
			<div class="alert alert-error mt-4">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="stroke-current shrink-0 h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
					/></svg
				>
				<span>{error.message}</span>
			</div>
		{/if}
	</main>

	<footer
		class="p-4 bg-base-100 border-t border-base-300 fixed bottom-0 w-full max-w-4xl mx-auto"
	>
		<form onsubmit={handleSubmit} class="flex gap-2">
			<input
				type="text"
				bind:value={input}
				placeholder="Type a message..."
				class="input input-bordered flex-1"
				disabled={isLoading}
			/>
			<button
				type="submit"
				class="btn btn-primary"
				disabled={isLoading || !input.trim()}
			>
				Send
			</button>
		</form>
	</footer>
</div>
