import { error } from '@sveltejs/kit';

/**
 * @type {import('@sveltejs/kit').PageLoad}
 */
export async function load({ fetch, params }) {
	const url = `http://localhost:8080/v1/bookmarks/`+ params.id;
	const res = await fetch(url, {
        credentials: 'include'
    });

	if (res.ok) {
		return {
			detail: await res.json()
		};
	}

	throw error(500, `Could not load ${url}`);
}
