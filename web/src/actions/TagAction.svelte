<script context="module" lang="ts">

import {baseURL} from "./../config.dev.js";

export async function GetTags() {

	const response = await fetch(baseURL + "/v1/tags", {
		method: 'GET'
	})

	if (response.ok) {
		return response.json();
	} else {
		throw new Error(response.statusText);
	}
}


export async function AddTagToBookmark(bookmarkID: number, tagName: string) {

	await fetch(baseURL + "/v1/tags", {
		method: 'POST',
		body: JSON.stringify({bookmarkID, tagName})
	}).then(response => {
		if(response.ok) {
			return response.json();
		}
	}).catch(error => {
		return error
	});
}


export async function DeleteTagFromBookmark(bookmarkID: number, tagName: string) {

	await fetch(baseURL + "/v1/tags", {
		method: 'DELETE',
		body: JSON.stringify({bookmarkID, tagName})
	}).then(response => {
		if(response.ok) {
			return response.json();
		}
	}).catch(error => {
		return error
	});
}

</script>