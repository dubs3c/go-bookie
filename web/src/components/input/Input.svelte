<script>
    import './input.css';

    export let placeholder = ""
    export let name = ""
    export let readOnly = false


</script>

<input 
    class="{['input']}"
    type="text"
    name="{name}"
    readOnly={readOnly}
    placeholder="{placeholder}"
/>