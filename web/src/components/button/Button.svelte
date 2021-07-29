<script>
    import './button.css';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    export let value = "Submit"
    export let large = false
    export let fullWidth = false
    export let danger = false

	function handleClick() {
		dispatch('buttonClick', {});
	}

</script>

<button
    on:click={handleClick}
    class="btn"
    class:large={large}
    class:fullWidth={fullWidth}
    class:danger={danger}>
    {value}</button>