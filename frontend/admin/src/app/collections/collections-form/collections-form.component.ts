import { Component, OnInit, ViewChild } from '@angular/core';
import { FormsModule, FormGroup } from "@angular/forms";
import { Collection } from "../../model/collection";
import { CollectionService } from "../../model/collection.service";

@Component({
  selector: 'app-collections-form',
  templateUrl: './collections-form.component.html',
  styleUrls: ['./collections-form.component.css']
})
export class CollectionsFormComponent implements OnInit {

  collection: Collection = new Collection();

  @ViewChild("collectionForm") collectionForm: FormGroup;

  constructor(private collectionService: CollectionService) { }

  ngOnInit() { 
  }

  onNameChange(data) {
    this.updateSlug(data);
  }

  onNameBlur(data) {
    
  }

  updateSlug(input) {
    this.collection.slug = input.trim().replace(/[^a-zA-Z0-9_]+/g, "-");
  }

  onSubmit() {
    console.log("onSubmit()")

    this.collectionService.save(this.collection);
  }

  get diagnostic() {
    return JSON.stringify(this.collection);
  }
}