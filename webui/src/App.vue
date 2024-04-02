<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
</script>

<template>

	<div class="container-fluid">
		<div class="row">
			<div class="col p-0">
				<main>
					<Sidebar :currentRoute="$route.name"/>
					<RouterView/>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
