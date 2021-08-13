<script lang="ts">
import { createEventDispatcher } from "svelte";

import "./tag.css"

export let value: string = ""
export let deletable: boolean = false
export let active: boolean = false
export let selectable: boolean = true

const dispatch = createEventDispatcher();

function handleTagSelection(name: string) {
    if(selectable) {
        active = !active
        dispatch('onActiveTag', {
            name: name,
        });
    }
}

function onDelete() {
    dispatch('onDeleteTag', {
        name: value,
    });
}

</script>
<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.1/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
<span on:click={() => handleTagSelection(value)} class="{(active === true && selectable == true) ? "tag highlight" : "tag"}">
    {value}
    
    {#if deletable}
        <i on:click="{onDelete}" class="close" title="Delete"><i class="fas fa-trash-alt"></i></i>
    {/if}
</span>