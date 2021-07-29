<script context="module">

import {baseURL} from "./../config.dev.ts";

export async function CreateBookmark() {
	const response = await fetch(baseURL, {
		method: 'POST',
		body: JSON.stringify({url})
	})

	if (response.ok) {
		return response.json();
	} else {
		throw new Error(response.statusText);
	}
}

</script>