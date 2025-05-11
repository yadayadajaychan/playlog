<script>
	import { onMount } from 'svelte';

	let playlog = $state(null);
	let loading = $state(true);
	let error = $state(null);

	let ascending = $state(false);
	let page = $state(1);
	let count = $state(50);
	$effect(() => {
		getPlaylog(ascending, page, count);
	});

	let selectedDate = $state(null);

	// onMount(async () => getPlaylog(ascending, page, count));

	async function getPlaylog(ascending, page, count) {
		try {
			const res = await fetch(`/api/playlog?ascending=${ascending}&page=${page}&count=${count}`);
			if (!res.ok) throw new Error('Failed to fetch data');
			playlog = await res.json();
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

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

	function scoreDiff(score1, score2) {
		if (score2 >= score1)
			return "+" + ((score2-score1)/10000);
		else
			return "-" + ((score1-score2)/10000);
	}
</script>

{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else}

<div id="top" class="w-screen max-w-128 flex flex-col">
	<div class="pt-4 pl-4 pr-4 pb-1 place-self-start">
		Source code for this website can be found at
		<a class="underline text-blue-500 whitespace-nowrap" href="https://github.com/yadayadajaychan/playlog">https://github.com/yadayadajaychan/playlog</a>
	</div>
	<div class="pt-4 pl-4 pr-4 pb-1 place-self-start">
		<label for="page">Page:</label>
		<input class="bg-gray-200" id="page" type="number" bind:value={page} min="1" max={playlog.MaxPage} />
		/ {playlog.MaxPage}
	</div>
	<div class="pt-1 pl-4 pr-4 pb-1 place-self-start">
		<label for="count">Count:</label>
		<input class="bg-gray-200" id="count" type="number" bind:value={count} min="1" />
	</div>
	<div class="pt-1 pl-4 pr-4 pb-1 place-self-start">
		<label for="ascending">Ascending:</label>
		<input id="ascending" type="checkbox" bind:checked={ascending}/>
	</div>

	{#each playlog.Playlog as entry}

	<div class="flex flex-col">

		<div class="flex flex-row p-4 space-x-4">
			<img class="self-start" width="75" height="75" src="/jacket/{entry.PlayInfo.SongId}.png" loading="lazy" alt="song jacket">
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
				<div class="whitespace-nowrap">
					{difficultyToString(entry.PlayInfo.Difficulty)}

					{#each entry.SongInfo.Charts as chart}
						{#if chart.Difficulty === entry.PlayInfo.Difficulty}
							{chart.InternalLevel / 10}
						{/if}
					{/each}
				</div>

				<div class="whitespace-nowrap">
					<strong>
					{entry.PlayInfo.Score/10000}%

					{scoreToRank(entry.PlayInfo.Score)}
					</strong>
				</div>

				{formatDate(entry.PlayInfo.UserPlayDate)}
			</div>
		</div>

		{#if selectedDate === entry.PlayInfo.UserPlayDate}
		<div class="flex flex-col p-4 space-y-2">
			<table class="table-auto text-right text-sm">
			<thead>
				<tr>
					<th class="p-1 border border-gray-400"></th>
					<th class="p-1 border border-gray-400">Critical<br>Perfect</th>
					<th class="p-1 border border-gray-400">Perfect</th>
					<th class="p-1 border border-gray-400">Great</th>
					<th class="p-1 border border-gray-400">Good</th>
					<th class="p-1 border border-gray-400">Miss</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<th class="p-1 border border-gray-400">Tap</th>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TapCriticalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TapPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TapGreat}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TapGood}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TapMiss}</td>
				</tr>
				<tr>
					<th class="p-1 border border-gray-400">Hold</th>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.HoldCriticalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.HoldPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.HoldGreat}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.HoldGood}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.HoldMiss}</td>
				</tr>
				<tr>
					<th class="p-1 border border-gray-400">Slide</th>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.SlideCriticalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.SlidePerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.SlideGreat}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.SlideGood}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.SlideMiss}</td>
				</tr>
				<tr>
					<th class="p-1 border border-gray-400">Touch</th>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TouchCriticalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TouchPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TouchGreat}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TouchGood}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TouchMiss}</td>
				</tr>
				<tr>
					<th class="p-1 border border-gray-400">Break</th>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.BreakCriticalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.BreakPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.BreakGreat}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.BreakGood}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.BreakMiss}</td>
				</tr>
				<tr>
					<th class="p-1 border border-gray-400">Total</th>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TotalCriticalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TotalPerfect}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TotalGreat}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TotalGood}</td>
					<td class="p-1 border border-gray-400">{entry.PlayInfo.TotalMiss}</td>
				</tr>
			</tbody>
			</table>

			<div>
				<strong>Fast:</strong> {entry.PlayInfo.FastCount} |
				<strong>Late:</strong> {entry.PlayInfo.LateCount}<br>
				<strong>Max Combo:</strong> {entry.PlayInfo.MaxCombo}/{entry.PlayInfo.TotalCombo}<br>
				<strong>Rating:</strong> {entry.PlayInfo.AfterRating} (+{entry.PlayInfo.AfterRating-entry.PlayInfo.BeforeRating})<br>
				<strong>Prev. Best:</strong> {entry.PreviousBestScore/10000}% ({scoreDiff(entry.PreviousBestScore, entry.PlayInfo.Score)}%)<br>
				<strong>Played with:</strong>
				{#each entry.PlayInfo.MatchingUsers as user}
					<div class="whitespace-nowrap">{user}</div>
				{/each}
			</div>
		</div>
		{/if}

	</div>

	{/each}

	<a class="p-4" href="#top">Back to Top</a>
</div>
{/if}
