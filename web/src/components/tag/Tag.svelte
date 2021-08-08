<script lang="ts">
import { createEventDispatcher } from "svelte";

import "./tag.css"

export let value = ""
export let deletable = false
export let active = false

const dispatch = createEventDispatcher();

function handleTagSelection(name: string) {
    active = !active
    dispatch('onActiveTag', {
        name: name,
    });
}

</script>

<span on:click={() => handleTagSelection(value)} class="{active === true ? "tag highlight" : "tag"}">
    {value}
    
    {#if deletable}
        <i class="close" title="Delete">x</i>
    {/if}
</span>