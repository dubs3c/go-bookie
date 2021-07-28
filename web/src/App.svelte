<script lang="ts">
import BookmarkScreen from "./components/bookmark/BookmarkScreen.svelte";
import Button from "./components/button/Button.svelte";
import Card from "./components/card/Card.svelte";
import Input from "./components/input/Input.svelte";
import Tag from "./components/tag/Tag.svelte";



</script>

<main>
	<div class="row">
		<div class="col">
			<h1>Go-Bookie!</h1>
			<hr />
		</div>
	</div>

	<Card>


		<div class="row">
			<div class="col">
				<h4>Filter</h4>
				<ul class="filter-list">
					<li><a href="?filter=unread"><i class="fas fa-angle-right"></i> Unread</a></li>
					<li><a href="?filter=read"><i class="fas fa-angle-right"></i> Read</a></li>
				</ul>
			</div>
			<div class="col lol">
				<h4>Save a link!</h4>
				<Input placeholder="Enter URL or some text"/>
			</div>
		</div>

		<div class="row">
			<div class="col">
				<h4>Tags</h4>
				<Tag value="APT"/>
				<Tag value="Spring Boot"/>
				<Tag value="Golang"/>
			</div>
		</div>

		<br />

		<Button value="Filter on tag"/>


	</Card>
	
	<br />
	
	<div class="row">
		<div class="col">
			<BookmarkScreen/>
		</div>
	</div>

</main>

<style>

	.filter-list li {
		list-style: none;
	}

	.lol {
		width: 80%;
	}

</style>