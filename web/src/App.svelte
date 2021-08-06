<script lang="ts">
import PureBookmarkList from "./components/bookmark/PureBookmarkList.svelte";
import Button from "./components/button/Button.svelte";
import Card from "./components/card/Card.svelte";
import Input from "./components/input/Input.svelte";
import Tag from "./components/tag/Tag.svelte";
import { DeleteBookmark, ArchiveBookmark, CreateBookmark, GetBookmarks, GetFilteredBookmarks } from "./actions/BookmarkAction.svelte";
import { onMount } from "svelte";
import { bookmarkStore, tagStore } from "./store";
import type { Pagination } from "./types/Pagination";
import type { TagType } from "./types/Tag";
import { GetTags } from "./actions/TagAction.svelte";

let url: string = "";
let promise = Promise.resolve([]);

let currentPage: number = 0
let totalPages: number = 0
let pageSize: number = 0

let archiveCheckbox: boolean = false
let deletedCheckbox: boolean = false

let activeTags: Record<string, string> = {}

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
	let paginatedObject: Pagination 
	paginatedObject = await GetFilteredBookmarks(currentPage, deletedCheckbox, archiveCheckbox, activeTags)
	
	currentPage = paginatedObject.page
	totalPages = paginatedObject.totalPages
	pageSize = paginatedObject.limit

	$bookmarkStore = paginatedObject.data
}

async function changePage(pageNumber: number) {
	currentPage = pageNumber
	let paginatedObject: Pagination = await GetBookmarks(currentPage)
	totalPages = paginatedObject.totalPages
	pageSize = paginatedObject.limit
	$bookmarkStore = paginatedObject.data
}

onMount(async () => {
	let paginatedObject: Pagination = await GetBookmarks(currentPage)
	currentPage = paginatedObject.page
	totalPages = paginatedObject.totalPages
	pageSize = paginatedObject.limit
	$bookmarkStore = paginatedObject.data

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
	if(activeTags[event.detail.name] === "") {
		delete activeTags[event.detail.name]
	} else {
		activeTags[event.detail.name] = ""
	}

	filterBookmarks()
}

</script>

<main>
	<div class="row">
		<div class="col">
			<h1><a href="/" style="text-decoration:none; color:turquoise;">SAMLA</a></h1>
			<hr />
		</div>
	</div>
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
							<p style="color: red">Could not create bookmark 😭 <strong>Error: {error.message}</strong></p>
						{/await}
					</div>
				</div>

				<div class="row">
					<div class="col">
						<h4>Filter</h4>
							<label><input type="checkbox" bind:checked={archiveCheckbox} on:change="{filterBookmarks}" /> Archived</label><br />
							<label><input type="checkbox" bind:checked={deletedCheckbox} on:change="{filterBookmarks}"/> Trash</label>
					</div>
				</div>
		
				<div class="row">
					<div class="col">
						<h4>Tags</h4>
						{#each $tagStore as tag}
							<Tag on:onActiveTag={onActiveTag} value="{tag.name}"/>
						{/each}
					</div>
				</div>
		
			</Card>
		</div>
		<div class="col">
			<PureBookmarkList
			bookmarks={$bookmarkStore}
			on:onDeleteBookmark={onDeleteBookmark}
			on:ArchiveTask={onArchiveTask}
			/>
			<div class="row">
				<p>
					{#if currentPage > 1 }
						<button on:click={() => changePage(1)}>« first</button> <button on:click={() => changePage(currentPage-1)}>previous</button>
					{/if}
					Page {currentPage} of {totalPages}
					{#if totalPages > 1 }
						<button on:click={() => changePage(currentPage+1)}>next</button>
						<button on:click={() => changePage(totalPages)}>last »</button>
					{/if}
				</p>
			</div>

		</div>
	</div>

	<br />
	<br />
</main>