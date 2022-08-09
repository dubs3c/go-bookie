<script lang="ts">
    import './input.css';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    export let placeholder: string = ""
    export let name: string = ""
    export let readOnly: boolean = false

    let inputText: string = ""

    function handleInput() {
        dispatch('inputText', {
            text: inputText
        });
    }

</script>

<input 
    bind:value={inputText}
    on:input={handleInput}
    class="{['input']}"
    type="password"
    name="{name}"
    readOnly={readOnly}
    placeholder="{placeholder}"
/>