
<script context="module">
/**
 * @type {import('@sveltejs/kit').Load}
 */
export async function load({}) {
    if(authStore.getToken() != "") {
        goto("/dashboard", {})
    }
}
</script>

<script lang="ts">
import Button from "../components/button/Button.svelte";
import Input from "../components/input/Input.svelte";
import Card from "../components/card/Card.svelte";
import PasswordInput from "../components/input/PasswordInput.svelte";
import { UserLogin } from "../actions/LoginAction.svelte";
import type { TokenResponse } from "../types/User";
import { authStore } from "../store";
import { goto } from '$app/navigation';

let email: string = ""
let password: string = ""

let promise = new Promise<TokenResponse>((resolve) => { resolve(<TokenResponse>{})});

function handleEmailInput(event) {
    email = event.detail.text
}

function handlePasswordInput(event) {
    password = event.detail.text
}

async function onSubmit() {
    promise = UserLogin(email, password)
    let val = await promise
    authStore.addToken(val.token)
    await goto("/dashboard", {})
}

</script>

<svelte:head>
	<title>SAMLA App | Login</title>
</svelte:head>

<div class="row">
    <div class="col">
        <h1>SAMLA | <span style="font-weight: normal">Login</span></h1>
    </div>
</div>

<div class="row">
    <div class="col">
        <Card>
            <Input on:inputText={ handleEmailInput }  placeholder="Email" />
            <PasswordInput on:inputText={handlePasswordInput} placeholder="Password" />

            <Button on:buttonClick="{onSubmit}" fullWidth={true} value="Login"/>
        </Card>
    </div>
    <div class="col">
        <p><strong>🚀 Welcome to the SAMLA App!</strong></p>
        <p>Login to view, add or read your bookmarks 😊 If you encounter any problems don't hesitate to contact support.</p>
        {#await promise}
            <p><i>Logging in...</i></p>
        {:catch error}
            <p style="color: red"><small>Could not login 😭 <strong>Error: {error.message}</strong></small></p>
        {/await}
    </div>

</div>