<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.1/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
<script>
    import "./bookmark.css"

    export let url = ""
    export let title = ""

</script>

<div class="bookmark">
    <h4><a href="{url}" target="_blank" rel="noopener noreferrer">{title}</a>
    </h4>
    <small>{url}</small>
    <div class="actions">
        <a href="\#" class="view">
            <i style=" cursor: pointer;" class="fas fa-eye"></i> View
        </a>
        <a href="\#" class="delete">
            <i style=" cursor: pointer;" class="fas fa-trash-alt"></i> Delete
        </a>
        <a href="\#" class="archive">
            <i style=" cursor: pointer;" class="fas fa-bookmark"></i> Archive
        </a>
    </div>
</div>
<div class="border"></div>