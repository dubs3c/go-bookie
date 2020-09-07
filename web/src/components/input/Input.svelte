<script>
    import './input.css';

    export let placeholder = ""
    export let name = ""
    export let readOnly = false
    export let buttonValue = "Submit"


</script>

<input 
    class="{['input']}"
    type="text"
    name="{name}"
    readOnly={readOnly}
    placeholder="{placeholder}"
/><input class="{['btn, btn-dark']}" type="button" value="{buttonValue}" />