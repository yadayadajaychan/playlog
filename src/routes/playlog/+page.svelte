<script>
	import { onMount } from 'svelte';

	let playlog = $state(null);
	let loading = $state(true);
	let error = $state(null);

	let selectedDate = $state(null);

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

	function selectDate(date) {
		selectedDate = (date === selectedDate) ? null : date;
	}
</script>

{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else}

<div class="w-screen max-w-128 flex flex-col">
	{#each playlog.Playlog as entry}

	<div class="flex flex-col">

		<div class="flex flex-row p-4 space-x-4">
			<img class="self-start" width="75" height="75" src="https://nijika.org/maimai/{entry.PlayInfo.SongId}.png" loading="lazy" alt="song jacket">
			<div class="w-1/2">
				<strong>{entry.SongInfo.Name}</strong><br>
				{entry.SongInfo.Artist}<br>
				<button class="underline text-blue-500" onclick={() => selectDate(entry.PlayInfo.UserPlayDate)}>
					{#if selectedDate === entry.PlayInfo.UserPlayDate}
						Hide
					{:else}
						Show
					{/if}
					Details
				</button>
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

		{#if selectedDate === entry.PlayInfo.UserPlayDate}
		<div>
			<table class="table-auto text-right">
			<thead>
				<tr>
					<th class="p-2 border border-gray-400"></th>
					<th class="p-2 border border-gray-400">Critical<br>Perfect</th>
					<th class="p-2 border border-gray-400">Perfect</th>
					<th class="p-2 border border-gray-400">Great</th>
					<th class="p-2 border border-gray-400">Good</th>
					<th class="p-2 border border-gray-400">Miss</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<th class="p-2 border border-gray-400">Tap</th>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TapCriticalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TapPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TapGreat}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TapGood}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TapMiss}</td>
				</tr>
				<tr>
					<th class="p-2 border border-gray-400">Hold</th>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.HoldCriticalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.HoldPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.HoldGreat}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.HoldGood}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.HoldMiss}</td>
				</tr>
				<tr>
					<th class="p-2 border border-gray-400">Slide</th>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.SlideCriticalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.SlidePerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.SlideGreat}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.SlideGood}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.SlideMiss}</td>
				</tr>
				<tr>
					<th class="p-2 border border-gray-400">Touch</th>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TouchCriticalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TouchPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TouchGreat}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TouchGood}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TouchMiss}</td>
				</tr>
				<tr>
					<th class="p-2 border border-gray-400">Break</th>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.BreakCriticalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.BreakPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.BreakGreat}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.BreakGood}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.BreakMiss}</td>
				</tr>
				<tr>
					<th class="p-2 border border-gray-400">Total</th>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TotalCriticalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TotalPerfect}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TotalGreat}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TotalGood}</td>
					<td class="p-2 border border-gray-400">{entry.PlayInfo.TotalMiss}</td>
				</tr>
			</tbody>
			</table>
		</div>
		{/if}

	</div>

	{/each}
</div>
{/if}
