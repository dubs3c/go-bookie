<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch, session, context }) {
		const url = `http://127.0.0.1:8080/v1/bookmarks/`+ page.params.id;
		const res = await fetch(url);

		if (res.ok) {
			return {
				props: {
					detail: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load ${url}`)
		};
	}
</script>

<script>
    import { page } from '$app/stores';
    import Card from '../../components/card/Card.svelte';
import Tag from '../../components/tag/Tag.svelte';

    export let detail;

</script>
<Card>
    <div class="row">
            <div class="col">
                <p><small><p><strong>Title:</strong></p></small></p>
                <h3>{detail.title}</h3>

                {#if detail.description }
                    <p><small><strong>Description:</strong></small></p>
                    <p>{detail.description}</p>
                {/if}

                {#if detail.image }
                    <p><small><strong>Image:</strong></small></p>
                    <img src="{detail.image}" width="100%" alt=""/>
                {/if}
            </div>
            
            <div class="col">
                <p><small><strong>URL:</strong> <a href="{detail.url}" rel="noopener noreferrer">{detail.url}</a></small></p>
                <p><small><strong>Saved:</strong> {detail.created_at}</small></p>
                <p><small><strong>Archived:</strong> {detail.archived}</small></p>
                <p><small><strong>Deleted:</strong> {detail.deleted}</small></p>
                {#if detail.tags }
                    <p><small><strong>Tags:</strong></small></p>
                    {#each detail.tags.split(",") as tag }
                        <Tag value="{tag}" deletable=true/>
                    {/each}
                {/if}
            </div>
    </div>
</Card>
<br />
<br />
<iframe src="http://127.0.0.1:8080/v1/bookmarks/{detail.id}?htmlbody=true" sandbox title="Body" width="100%" height="600px">
    <p>Your browser does not support iframes.</p>
</iframe>

<br />