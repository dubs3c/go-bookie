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

export async function GetBookmarks() {
	const response = await fetch(baseURL + "/api/v1/bookmarks", {
		method: 'GET'
	})

	if (response.ok) {
		return response.json();
	} else {
		throw new Error(response.statusText);
	}
}

</script>