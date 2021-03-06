<template>
  <v-dialog
      :value="value"
      @input="emitChange"
      max-width="400px"
  >
    <v-card>
      <v-card-title>
        <span class="headline">Update a new Tag</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-form ref="form" v-model="valid">
            <v-row>
              <v-text-field
                  label="Name"
                  v-model="tag.name"
                  @keydown.enter="saveTag"
                  :rules="rules.name"
              ></v-text-field>
            </v-row>
            <v-row>
              <v-text-field
                  label="Description"
                  v-model="tag.description"
                  @keydown.enter="saveTag"
                  required
              ></v-text-field>
            </v-row>
            <v-row v-if="hasAdminRole">
              <v-checkbox
                  label="Basic"
                  v-model="tag.basic"
                  disabled
              ></v-checkbox>
            </v-row>
            <v-row v-if="!tag.basic">
              <v-select
                  label="Parent"
                  v-model="tag.parent"
                  item-text="name"
                  item-value="id"
                  :items="tags"
                  clearable
              ></v-select>
            </v-row>
            <v-row v-if="!tag.basic">
              <v-select
                  label="Project"
                  v-model="tag.project"
                  item-text="name"
                  item-value="id"
                  :items="projects"
                  clearable
              ></v-select>
            </v-row>
            <v-row>
              <tag-color-picker v-model="tag.color"/>
            </v-row>
          </v-form>
        </v-container>
      </v-card-text>

      <div class="d-flex justify-center">
        <div class="pl-3 pr-3 global-form-error error--text">
          {{ globalFormError }}
        </div>
      </div>

      <v-card-actions>
        <v-btn
            color="blue darken-1"
            text
            @click="saveTag"
        >
          Save
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="blue darken-1"
            text
            @click="emitChange(false)"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import TagService from "@/service/TagService"
import TagColorPicker from "@/components/tags/TagColorPicker"

export default {
  name: "CreateTagDialog",
  components: {
    TagColorPicker,
  },

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    tag: {
      type: Object,
      required: true,
    },
    tags: {
      type: Array,
      required: true,
    },
    projects: {
      type: Array,
      required: true,
    },
    hasAdminRole: {
      type: Boolean,
      required: true,
    },
  },
  data: () => ({
    valid: true,
    globalFormError: '',
    rules: {
      name: [
        v => !!v || 'Name is required',
      ],
    },
  }),
  methods: {
    saveTag() {
      if (this.$refs.form.validate()) {
        TagService.update(this.tag)
            .then(() => {
              this.$emit('refreshTags', null)
              this.emitChange(false)
            })
            .catch(error => {
              if (error && error.response && (error.response.status === 400 || error.response.status === 404)) {
                this.globalFormError = error.response.data.error
              }
            })
      }
    },
    emitChange(value) {
      this.$emit('input', value)
    }
  },
  watch: {
    'value': function () {
      if (this.value === false) {
        this.globalFormError = ''
      }
    }
  },
}
</script>