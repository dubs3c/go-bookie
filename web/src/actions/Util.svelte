
<script context="module" lang="ts">
import { authStore } from "../store";

export async function AuthenticatedFetch(url: string, stuff?: RequestInit): Promise<Response> {
    
    const token: string = authStore.getToken()

    if(token === "") {
        throw new Error("you aint logged in, fool")
    }

    stuff.headers = {
        "Authorization": "Bearer " + token
    }

    return await fetch(url, stuff)
}
    
</script>