<script lang="ts">
	import { Chat } from '@ai-sdk/svelte';

	const chat = new Chat({
		api: '/chat',
		fetch: async (input, init) => {
			const res = await fetch(input, init);
			if (!res.ok) return res;

			// Adapt Linden's strict SSE stream into Vercel AI SDK's TextStream format
			const transform = new TransformStream({
				transform(chunk, controller) {
					const text = new TextDecoder().decode(chunk);
					const lines = text.split('\n');
					for (let line of lines) {
						if (line.startsWith('data: ')) {
							try {
								const data = JSON.parse(line.slice(6));
								if (data.text) {
									// Vercel AI text chunk format is "0:\"chunk\""
									controller.enqueue(
										new TextEncoder().encode(`0:${JSON.stringify(data.text)}\n`)
									);
								}
							} catch (e) {
								// skip invalid json
							}
						}
					}
				}
			});

			return new Response(res.body?.pipeThrough(transform), {
				headers: { 'Content-Type': 'text/plain; charset=utf-8' }
			});
		}
	});
</script>

<div class="flex flex-col h-screen bg-base-100 text-base-content max-w-4xl mx-auto">
	<header class="p-4 border-b border-base-300">
		<h1 class="text-xl font-bold">Linden AI</h1>
	</header>

	<main class="flex-1 overflow-y-auto p-4 space-y-4 pb-32">
		{#each chat.messages as msg}
			<div class="chat {msg.role === 'user' ? 'chat-end' : 'chat-start'}">
				<div class="chat-header capitalize mb-1 opacity-50">
					{msg.role}
				</div>
				<div class="chat-bubble {msg.role === 'user' ? 'chat-bubble-primary' : 'chat-bubble-neutral'} whitespace-pre-wrap">
					{msg.content}
				</div>
			</div>
		{/each}

		{#if chat.isLoading && chat.messages[chat.messages.length - 1]?.role === 'user'}
			<div class="chat chat-start">
				<div class="chat-bubble chat-bubble-neutral opacity-50">
					<span class="loading loading-dots loading-sm"></span>
				</div>
			</div>
		{/if}

		{#if chat.error}
			<div class="alert alert-error mt-4">
				<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
				<span>{chat.error.message}</span>
			</div>
		{/if}
	</main>

	<footer class="p-4 bg-base-100 border-t border-base-300 fixed bottom-0 w-full max-w-4xl mx-auto">
		<form onsubmit={chat.handleSubmit} class="flex gap-2">
			<input
				type="text"
				bind:value={chat.input}
				placeholder="Type a message..."
				class="input input-bordered flex-1"
				disabled={chat.isLoading}
			/>
			<button type="submit" class="btn btn-primary" disabled={chat.isLoading || !chat.input?.trim()}>
				Send
			</button>
		</form>
	</footer>
</div>
