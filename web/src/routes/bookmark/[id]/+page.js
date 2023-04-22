import { error } from '@sveltejs/kit';

/** @type {import('./$types').PageLoad} */
export async function load({ fetch, params }) {
	const url = `http://localhost:8181/v1/bookmarks/`+ params.id;
    console.log(url);
	const res = await fetch(url, {
        credentials: 'include'
    });

	if (res.ok) {
		return {
			data: await res.json()
		};
	} else {
        console.log(res.status);
        throw error(500, `Could not load ${url}`);
    }
}
