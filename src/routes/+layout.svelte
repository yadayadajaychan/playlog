<script lang="ts">
	import '../app.css';
	
	let { children } = $props();
</script>

<nav>
<ul class="pt-4 pl-4 pr-4 space-x-2 flex flex-row text-blue-500 underline">
	<li>
		<a href="/">Playlog</a>
	</li>
	<li>
		<a href="/about/">About</a>
	</li>
</ul>
</nav>

{@render children()}
