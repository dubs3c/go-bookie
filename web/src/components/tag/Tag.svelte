<script>
    import "./tag.css"
    export let value = ""
    export let deletable = false
</script>

<div class="tag-container">
    <div class="tag">
        {value}

        {#if deletable}
        <i class="close" title="Delete">x</i>
        {/if}
    </div>
</div>