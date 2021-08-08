<script>
    import './button.css';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    export let value = "Submit"
    export let large = false
    export let small = false
    export let fullWidth = false
    export let danger = false
    export let success = false

	function handleClick() {
		dispatch('buttonClick', {});
	}

</script>

<button
    on:click={handleClick}
    class="btn"
    class:large={large}
    class:small={small}
    class:fullWidth={fullWidth}
    class:danger={danger}
    class:success={success}>
    {value}</button>