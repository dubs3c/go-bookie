<script lang="ts">
import PureBookmarkList from "./../components/bookmark/PureBookmarkList.svelte";
import Button from "./../components/button/Button.svelte";
import Card from "./../components/card/Card.svelte";
import Input from "./../components/input/Input.svelte";
import Tag from "./../components/tag/Tag.svelte";
import { CreateBookmark, GetBookmarks, GetFilteredBookmarks } from "./../actions/BookmarkAction.svelte";
import { onMount } from "svelte";
import { bookmarkStore, settingsStore, tagStore } from "./../store";
import type { Pagination } from "./../types/Pagination";
import type { TagType } from "./../types/Tag";
import { GetTags } from "./../actions/TagAction.svelte";


let url: string = "";
let promise = Promise.resolve([]);

function handleBookmarkInput(event) {
	url = event.detail.text
}

function handleClick() {
	if(url != "") {
		promise = CreateBookmark(url)
		// run GetBookmarks() ?
	}
}

async function filterBookmarks() {
	// Reset current page to 1 to ensure when a new filter is made, we request page 1 first
	$settingsStore.currentPage = 1
	let paginatedObject: Pagination 
	paginatedObject = await GetFilteredBookmarks(
		$settingsStore.currentPage,
		$settingsStore.deletedChecked,
		$settingsStore.archiveChecked,
		$settingsStore.activeTags
	)

	$settingsStore.currentPage = paginatedObject.page
	$settingsStore.totalPages = paginatedObject.totalPages
	$settingsStore.pageSize = paginatedObject.limit

	$bookmarkStore = paginatedObject.data
}

async function changePage(pageNumber: number) {
	$settingsStore.currentPage = pageNumber
	let paginatedObject: Pagination = await GetBookmarks($settingsStore.currentPage)
	$settingsStore.totalPages = paginatedObject.totalPages
	$settingsStore.pageSize = paginatedObject.limit
	$settingsStore.currentPage = paginatedObject.page
	$bookmarkStore = paginatedObject.data
}

onMount(async () => {
	await filterBookmarks()

	let tags: TagType[] = await GetTags()
	$tagStore = tags
})

async function onDeleteBookmark(event) {
	filterBookmarks()
}

async function onArchiveTask(event) {
	filterBookmarks()
}

function onActiveTag(event) {
	let activeTags = $settingsStore.activeTags
	if(activeTags[event.detail.name] === "") {
		delete activeTags[event.detail.name]
	} else {
		activeTags[event.detail.name] = ""
	}

	filterBookmarks()
}

</script>

<svelte:head>
	<title>SAMLA App</title>
</svelte:head>


<div class="row">
	<div class="col">
		<Card>
			<div class="row">
				<div class="col">
					<h4>Save a link!</h4>
					<Input on:inputText={ handleBookmarkInput } placeholder="Enter URL or some text"/>
					<Button on:buttonClick={ handleClick } value="Save bookmark!"/>
					{#await promise}
						<p><i>Adding bookmark...</i></p>
					{:catch error}
						<p style="color: red"><small>Could not create bookmark 😭 <strong>Error: {error.message}</strong></small></p>
					{/await}
				</div>
			</div>

			<div class="row">
				<div class="col">
					<h4>Filter</h4>
						<label><input type="checkbox" bind:checked={$settingsStore.archiveChecked} on:change="{filterBookmarks}" /> Archived</label><br />
						<label><input type="checkbox" bind:checked={$settingsStore.deletedChecked} on:change="{filterBookmarks}"/> Trash</label>
				</div>
			</div>
	
			<div class="row">
				<div class="col">
					<h4>Tags</h4>
					{#each $tagStore as tag}
						<Tag on:onActiveTag={onActiveTag} value="{tag.name}" active="{$settingsStore.activeTags[tag.name] === "" ? true : false}" />
					{/each}
				</div>
			</div>
	
		</Card>
	</div>
	<div class="col">
		{#if $settingsStore.deletedChecked}
			<small style="color:lightsteelblue;"><i>Items in trash will be deleted after 30 days</i></small>
		{/if}
		<PureBookmarkList
		bookmarks={$bookmarkStore}
		on:onDeleteBookmark={onDeleteBookmark}
		on:ArchiveTask={onArchiveTask}
		/>
		<div class="row">
			<p>
				{#if $settingsStore.currentPage > 1 }
					<button on:click={() => changePage(1)}>« first</button> <button on:click={() => changePage($settingsStore.currentPage-1)}>previous</button>
				{/if}
				<!--
					TODO 1: Bug here regarding totalPages being incorrent, sent by the server				
				-->
				Page {$settingsStore.currentPage} of {$settingsStore.totalPages}
				{#if $settingsStore.totalPages > 1 }
					<button on:click={() => changePage($settingsStore.currentPage+1)}>next</button>
					<button on:click={() => changePage($settingsStore.totalPages)}>last »</button>
				{/if}
			</p>
		</div>

	</div>
</div>

<br />
<br />
