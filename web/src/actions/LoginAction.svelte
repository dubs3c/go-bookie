
<script context="module" lang="ts">

import {baseURL} from "./../config.dev.js";
import type { TokenResponse } from "../types/User";

export async function UserLogin(email: string, password: string): Promise<TokenResponse> {
    const response = await fetch(baseURL + "/v1/login", {
        method: 'POST',
        body: JSON.stringify({email, password})
    })

    if (response.ok) {
        return response.json() as Promise<TokenResponse>
    } else {
        let resp = await response.json()
        throw new Error(response.statusText + ": " + resp.error);
    }
}

</script>