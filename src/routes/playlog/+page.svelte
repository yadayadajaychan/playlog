<script>
	import { onMount } from 'svelte';

	let playlog = null;
	let loading = true;
	let error = null;


	onMount(async () => {
		try {
			const res = await fetch('/api/playlog?count=100');
			if (!res.ok) throw new Error('Failed to fetch data');
			playlog = await res.json();
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	});

	function difficultyToString(difficulty) {
		switch (difficulty) {
		case 0:
			return "BASIC"
		case 1:
			return "ADVANCED"
		case 2:
			return "EXPERT"
		case 3:
			return "MASTER"
		case 4:
			return "REMASTER"
		case 5:
			return "UTAGE"
		default:
			return "UNKNOWN"
		}
	}

	function scoreToRank(score) {
		if (score >= 1005000)
			return "SSS+"
		else if (score >= 1000000)
			return "SSS"
		else if (score >= 995000)
			return "SS+"
		else if (score >= 990000)
			return "SS"
		else if (score >= 980000)
			return "S+"
		else if (score >= 970000)
			return "S"
		else if (score >= 940000)
			return "AAA"
		else if (score >= 900000)
			return "AA"
		else if (score >= 800000)
			return "A"
		else
			return "X"
	}

	function formatDate(ts) {
		const date = new Date(ts * 1000);
		const pad = n => n.toString().padStart(2, '0');

		const year = date.getFullYear();
		const month = pad(date.getMonth() + 1);
		const day = pad(date.getDate());
		const hours = pad(date.getHours());
		const minutes = pad(date.getMinutes());

		return `${year}-${month}-${day} ${hours}:${minutes}`;
	}
</script>

{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else}
	<div class="w-screen max-w-128 flex flex-col">
	{#each playlog.Playlog as entry}
		<div class="flex flex-row p-4 space-x-4">
			<img class="self-start" width="75" height="75" src="https://nijika.org/maimai/{entry.PlayInfo.SongId}.png" loading="lazy" alt="song jacket">
			<div class="w-1/2">
				<strong>{entry.SongInfo.Name}</strong><br>
				{entry.SongInfo.Artist}
			</div>

			<div>
				{difficultyToString(entry.PlayInfo.Difficulty)}

				{#each entry.SongInfo.Charts as chart}
					{#if chart.Difficulty === entry.PlayInfo.Difficulty}
						{chart.InternalLevel / 10}
					{/if}
				{/each}
				<br>

				<strong>
				{entry.PlayInfo.Score/10000}%

				{scoreToRank(entry.PlayInfo.Score)}
				</strong>
				<br>

				{formatDate(entry.PlayInfo.UserPlayDate)}
			</div>
		</div>
	{/each}
	</div>
{/if}
