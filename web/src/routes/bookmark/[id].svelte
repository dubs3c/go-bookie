<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ fetch, params }) {
		const url = `http://localhost:8080/v1/bookmarks/`+ params.id;
		const res = await fetch(url);

		if (res.ok) {
			return {
				props: {
					detail: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load ${url}`)
		};
	}
</script>
<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.1/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
<script lang="ts">
import Card from '../../components/card/Card.svelte';
import Tag from '../../components/tag/Tag.svelte';
import { goto } from '$app/navigation';
import { AddTagToBookmark, DeleteTagFromBookmark } from '../../actions/TagAction.svelte';
import { onMount } from 'svelte';
import type { Bookmark } from '../../types/Bookmark';

export let detail: Bookmark;

let newTag: string =  ""

let ref: string = ""

$: currentTags = []

async function onKeyPress(event) {
    if(event.charCode === 13) {
        await AddTagToBookmark(detail.id, newTag)
        .then(() => {
            if(currentTags[0] === "") {
                currentTags = [newTag]
            } else {
                currentTags = [...currentTags, newTag]
            }
        }).catch(error => {
            console.log(error)
        });
        newTag = ""
    }
}

async function deleteTag(event) {
    let name: string = event.detail.name

    await DeleteTagFromBookmark(detail.id, name)

    // meh, svelte reactivity only triggers on assignment
    currentTags.splice(currentTags.indexOf(name), 1)
    currentTags = [...currentTags]
    console.log(currentTags)

}
onMount(() => {
    ref = document.referrer;
    currentTags = detail.tags.split(",")
})


</script>

<button class="back" title="Go back" on:click={() => goto(ref.length > 0 ? ref : "/")}><i class="fas fa-arrow-alt-circle-left"></i> Go back</button>

<Card>
    <div class="row">
            <div class="col">
                <p><small><p><strong>Title:</strong></p></small></p>
                <h3>{detail.title}</h3>

                {#if detail.description }
                    <p><small><strong>Description:</strong></small></p>
                    <p>{detail.description}</p>
                {/if}

                {#if currentTags[0] != "" }
                    <p><small><strong>Tags:</strong></small></p>
                    {#each currentTags as tag }
                        <Tag value="{tag}" deletable={true} on:onDeleteTag="{deleteTag}"/>
                    {/each}
                {/if}
                <input class="tag-input" placeholder="Press enter to add tag" bind:value="{newTag}" on:keypress="{onKeyPress}"/>
            </div>
            
            <div class="col">
                <p><small><strong>URL:</strong> <a href="{detail.url}" rel="noopener noreferrer">{detail.url}</a></small></p>
                <p><small><strong>Saved:</strong> {detail.createdAt}</small></p>
                <p><small><strong>Archived:</strong> {detail.archived}</small></p>
                <p><small><strong>Deleted:</strong> {detail.deleted}</small></p>
            </div>
    </div>
</Card>
<br />
<br />
{#if detail.body }
    <iframe src="http://localhost:8080/v1/bookmarks/{detail.id}?htmlbody=true" sandbox title="Body" width="100%" height="600px">
        <p>Your browser does not support iframes.</p>
    </iframe>

    <br />
{/if}


<style>
    .back {
        color: #212529;
        background-color: #f8f9fa;
        border-color: #f8f9fa;
        display: inline-block;
        font-weight: 400;
        padding: .375rem .75rem;
        font-size: 1rem;
        line-height: 1.5;
        border-radius: .25rem;
        transition: color .15s ease-in-out,background-color .15s ease-in-out,border-color .15s ease-in-out,box-shadow .15s ease-in-out;
        border: 1px solid transparent;
        margin-bottom: .5rem !important;
        cursor: pointer;
    }

    .tag-input {
        border: none;
        border-bottom: 1px solid #ccc;
        padding: 0.5em;
        outline: none;
        display: block;
        margin-top: 0.5em;
        font-size: medium;
    }

</style>