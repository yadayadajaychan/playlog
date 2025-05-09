<script>
	import { onMount } from 'svelte';

	let playlog = null;
	let loading = true;
	let error = null;

	onMount(async () => {
		try {
			const res = await fetch('/api/playlog');
			if (!res.ok) throw new Error('Failed to fetch data');
			playlog = await res.json();
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else}
	{#each playlog.Playlog as entry}
		<div>
		<img width="95" height="95" src="https://nijika.org/maimai/{entry.PlayInfo.SongId}.png" loading="lazy" style="display: inline;">
		<div>
			<strong>{entry.SongInfo.Name}</strong><br>
			{entry.SongInfo.Artist}
		</div>
		</div>
	{/each}
	<pre>{JSON.stringify(playlog, null, 2)}</pre>
{/if}
