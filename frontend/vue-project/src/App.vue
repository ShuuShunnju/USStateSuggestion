<template>
  <div id="app">
    <input v-model="searchQuery" @input="fetchStates" placeholder="Search for a state" />
    <ul>
      <li v-for="state in uniqueStates" :key="state.id" @click="selectState(state)">
        {{ state.name }}
      </li>
    </ul>
    <div ref="map" id="map"></div>
  </div>
</template>

<script>
import { gql } from '@apollo/client/core';
import { Loader } from '@googlemaps/js-api-loader';

const STATES_QUERY = gql`
  query States($search: String) {
    states(search: $search) {
      id
      name
    }
  }
`;

export default {
  data() {
    return {
      searchQuery: '',
      states: [],
      map: null,
      selectedState: null,
    };
  },
  computed: {
    uniqueStates() {
      const stateMap = new Map();
      this.states.forEach((state) => {
        if (!stateMap.has(state.id)) {
          stateMap.set(state.id, state);
        }
      });
      return Array.from(stateMap.values());
    },
  },
  methods: {
    fetchStates() {
      this.$apollo
        .query({
          query: STATES_QUERY,
          variables: { search: this.searchQuery },
        })
        .then((result) => {
          this.states = result.data.states;
        })
        .catch((error) => {
          console.error('Error fetching states:', error);
        });
    },
    selectState(state) {
      this.selectedState = state;
      this.highlightStateOnMap(state.name);
    },
    highlightStateOnMap(stateName) {
      if (!this.map) return;

      const geocoder = new google.maps.Geocoder();
      geocoder.geocode({ address: stateName }, (results, status) => {
        if (status === 'OK') {
          this.map.setCenter(results[0].geometry.location);
          new google.maps.Marker({
            map: this.map,
            position: results[0].geometry.location,
          });
        } else {
          console.error('Geocode was not successful for the following reason:', status);
        }
      });
    },
  },
  mounted() {
    const loader = new Loader({
      apiKey: 'AIzaSyDTXcJtPYQX3-pWgVSqoXUApS_p_zYoDwI',
      version: 'weekly',
    });

    loader
      .load()
      .then(() => {
        this.$nextTick(() => {
          if (this.$refs.map) {
            this.map = new google.maps.Map(this.$refs.map, {
              center: { lat: 37.0902, lng: -95.7129 },
              zoom: 4,
            });
          } else {
            console.error('Map element not found');
          }
        });
      })
      .catch((error) => {
        console.error('Error loading Google Maps:', error);
      });
  },
};
</script>

<style>
#map {
  height: 500px;
  width: 100%;
}
</style>
