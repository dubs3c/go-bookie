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

</script>