<script lang="ts">
import PureBookmarkList from "./components/bookmark/PureBookmarkList.svelte";
import Button from "./components/button/Button.svelte";
import Card from "./components/card/Card.svelte";
import Input from "./components/input/Input.svelte";
import Tag from "./components/tag/Tag.svelte";
import { CreateBookmark, GetBookmarks } from "./actions/BookmarkAction.svelte";
import { onMount } from "svelte";
import { bookmarkStore } from "./store";
import type { Pagination } from "./types/Pagination";

let url: string = "";
let promise = Promise.resolve([]);

let currentPage: number = 0
let totalPages: number = 0
let pageSize: number = 0

function handleBookmarkInput(event) {
	url = event.detail.text
}

function handleClick() {
	if(url != "") {
		promise = CreateBookmark(url)
		// run GetBookmarks() ?
	}
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


function onDeleteBookmark(event) {
  bookmarkStore.deleteBookmark(event.detail.id);
}

function onArchiveTask(event) {
	bookmarkStore.archiveBookmark(event.detail.id);
}

</script>

<main>
	<div class="row">
		<div class="col">
			<h1>Go-Bookie!</h1>
			<hr />
		</div>
	</div>

	<Card>
		<div class="row">
			<div class="col">
				<h4>Filter</h4>
				<ul class="filter-list">
					<li><a href="?filter=unread"><i class="fas fa-angle-right"></i> Unread</a></li>
					<li><a href="?filter=read"><i class="fas fa-angle-right"></i> Read</a></li>
				</ul>
			</div>
			<div class="col lol">
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
				<h4>Tags</h4>
				<Tag value="APT"/>
				<Tag value="Spring Boot"/>
				<Tag value="Golang"/>
			</div>
		</div>

		<br />

		<Button value="Filter on tag"/>

	</Card>
	
	<br />
	<div class="row">
		<div class="col">
			<PureBookmarkList
			bookmarks={$bookmarkStore}
			on:onDeleteBookmark={onDeleteBookmark}
			on:ArchiveTask={onArchiveTask}
		  />
		</div>
	</div>
	<div class="row">
		<p>
			{#if currentPage > 1 }
				<button on:click={() => changePage(1)}>« first</button> <button on:click={() => changePage(currentPage-1)}>previous</button>
			{/if}
			Page {currentPage} of {totalPages}
			<button on:click={() => changePage(currentPage+1)}>next</button>
			<button on:click={() => changePage(totalPages)}>last »</button>
		</p>
	</div>
	<br />
</main>

<style>
	.filter-list li {
		list-style: none;
	}

	.lol {
		width: 80%;
	}
</style>