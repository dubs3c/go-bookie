<script lang="ts">
import PureBookmarkList from "./components/bookmark/PureBookmarkList.svelte";
import Button from "./components/button/Button.svelte";
import Card from "./components/card/Card.svelte";
import Input from "./components/input/Input.svelte";
import Tag from "./components/tag/Tag.svelte";
import { DeleteBookmark, ArchiveBookmark, CreateBookmark, GetBookmarks, GetFilteredBookmarks } from "./actions/BookmarkAction.svelte";
import { onMount } from "svelte";
import { bookmarkStore } from "./store";
import type { Pagination } from "./types/Pagination";

let url: string = "";
let promise = Promise.resolve([]);

let currentPage: number = 0
let totalPages: number = 0
let pageSize: number = 0

let archiveCheckbox: boolean = false
let deletedCheckbox: boolean = false

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
	
	if(deletedCheckbox == false && archiveCheckbox == false) {
		paginatedObject= await GetBookmarks(currentPage)
	} else {
		paginatedObject = await GetFilteredBookmarks(currentPage, deletedCheckbox, archiveCheckbox)
	}

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
})

async function onDeleteBookmark(event) {
	bookmarkStore.deleteBookmark(event.detail.id);
	await DeleteBookmark(event.detail.id)
	const aa: Pagination = await GetBookmarks(currentPage)
	$bookmarkStore = aa.data
}

async function onArchiveTask(event) {
	bookmarkStore.archiveBookmark(event.detail.id);
	await ArchiveBookmark(event.detail.id)
	const aa: Pagination = await GetBookmarks(currentPage)
	$bookmarkStore = aa.data
}

</script>

<main>
	<div class="row">
		<div class="col">
			<h1>Go-Bookie!</h1>
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
							<label><input type="checkbox" bind:checked={deletedCheckbox} on:change="{filterBookmarks}"/> Deleted</label>
					</div>
				</div>
		
				<div class="row">
					<div class="col">
						<h4>Tags</h4>
						<Tag value="APT"/>
						<Tag value="Spring Boot"/>
						<Tag value="Golang"/>
					</div>
				</div>
		
				<br />
		
				<Button value="Filter on tag"/>
		
			</Card>
		</div>
		<div class="col">

			<PureBookmarkList
			bookmarks={$bookmarkStore}
			on:onDeleteBookmark={onDeleteBookmark}
			on:ArchiveTask={onArchiveTask}
			/>

		</div>
	</div>
	
	
	<br />

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
	<br />
</main>