import { Component, OnInit, ViewChild } from '@angular/core';
import { FormsModule, FormGroup } from "@angular/forms";
import { Collection } from "../../model/collection";

@Component({
  selector: 'app-collections-form',
  templateUrl: './collections-form.component.html',
  styleUrls: ['./collections-form.component.css']
})
export class CollectionsFormComponent implements OnInit {

  collection: Collection = new Collection();

  @ViewChild("collectionForm") collectionForm: FormGroup;

  constructor() { }

  ngOnInit() { }

  onNameChange(data) {
    this.collection.slug = data.replace(/[^a-zA-Z0-9_]+/g, "-");
  }

  onSubmit() {
    console.log("onSubmit()")
  }

  get diagnostic() {
    return JSON.stringify(this.collection);
  }
}