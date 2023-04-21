<script context="module" lang="ts">

    import {baseURL} from "./../config.dev.js";
    
    export async function Login(email, password) {

        const response = await fetch(baseURL + "/login", {
            method: 'POST',
            body: JSON.stringify({email, password}),
            credentials: 'include'
        })

        if (response.ok) {
            return true
        } else {
            throw new Error(response.statusText);
        }
    }

</script>