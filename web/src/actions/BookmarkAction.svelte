<script context="module" lang="ts">

import {baseURL} from "./../config.dev.js";

export async function CreateBookmark(url: string) {
	const response = await fetch(baseURL + "/api/v1/bookmarks", {
		method: 'POST',
		body: JSON.stringify({url})
	})

	if (response.ok) {
		return response.json();
	} else {
		throw new Error(response.statusText);
	}
}

export async function GetBookmarks(page: number) {
	if(page == 0) {
		page = 1
	}
	const response = await fetch(baseURL + "/api/v1/bookmarks?page=" + page, {
		method: 'GET'
	})

	if (response.ok) {
		return response.json();
	} else {
		throw new Error(response.statusText);
	}
}

export async function GetFilteredBookmarks(page: number, deleted: boolean, archived: boolean) {
	if(page == 0) {
		page = 1
	}

	let filter: string = ""

	if(archived) {
		filter += "&archived=true"
	}

	if(deleted) {
		filter += "&deleted=true"
	}

	let url: string = "/api/v1/bookmarks?page=" + page

	if(filter != "") {
		url += filter
	}

	const response = await fetch(baseURL + url, {
		method: 'GET'
	})

	if (response.ok) {
		return response.json();
	} else {
		throw new Error(response.statusText);
	}
}

export async function DeleteBookmark(id: number) {
	await fetch(baseURL + "/api/v1/bookmarks/"+id, {
		method: "DELETE"
	}).then(response => {
		if(response.ok) {
			return response.json();
		}
	}).catch(error => {
		console.log(error)
		return error
	})
}

export async function ArchiveBookmark(id: number) {
	let archived = true
	await fetch(baseURL + "/api/v1/bookmarks/"+id, {
		method: "PATCH",
		body: JSON.stringify({archived})
	}).then(response => {
		if(response.ok) {
			return response.json();
		}
	}).catch(error => {
		console.log(error)
		return error
	})
}

</script>