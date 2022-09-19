<script lang="ts">
import Button from "../../components/button/Button.svelte";
import Card from "../../components/card/Card.svelte";
import Input from "../../components/input/Input.svelte";
import InputPassword from "../../components/input/InputPassword.svelte";
import { Login } from '../../actions/LoginAction.svelte';
import { goto } from '$app/navigation';
import { authenticated } from "../../store";
import { onMount } from "svelte";

// TODO - Create middleware that makes sure the user is signed in, otherwise redirect to login

let username = "";
let password = "";

async function onLogin() {

    if(username === "" || password === "") {
        console.log("please enter username and password");
        return
    }

    await Login(username, password).then(() => {
        // store something in store, indicating we are logged in
        window.localStorage.setItem("authenticated", "true");
        goto("/dashboard")
    }).catch(error => {
        console.log("Login Error:")
        console.log(error)
    });
}

onMount(async () => {
    if(window.localStorage.getItem("authenticated") === "true") {
        goto("/dashboard");
    }
})

</script>

<svelte:head>
	<title>SAMLA App | Login </title>
</svelte:head>

<br />
<br />
<br />
<div class="row">
    <div class="col">
        <h3><strong>Login</strong></h3>
        <Card>
            <Input on:inputText={ (event) => {username = event.detail.text} } placeholder="Email" />
            <InputPassword on:inputText={ (event) => {password = event.detail.text} } placeholder="Password" />
            <Button on:buttonClick={ onLogin }  />
        </Card>
    </div>
</div>