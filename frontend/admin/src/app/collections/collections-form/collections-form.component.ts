import { Component, OnInit, ViewChild } from '@angular/core';
import { FormsModule, FormGroup } from "@angular/forms";
import { Router } from "@angular/router";
import { Collection } from "../../model/collection";
import { CollectionService } from "../../model/collection.service";

@Component({
  selector: 'app-collections-form',
  templateUrl: './collections-form.component.html',
  styleUrls: ['./collections-form.component.css']
})
export class CollectionsFormComponent implements OnInit {

  collection: Collection = new Collection();

  submitting: boolean = false;

  @ViewChild("collectionForm") collectionForm: FormGroup;

  constructor(
    private router: Router,
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
  }

  onNameChange(data) {
    this.updateSlug(data);
  }

  updateSlug(input) {
    this.collection.slug = input.trim().replace(/[^a-zA-Z0-9_]+/g, "-");
  }

  onSubmit() {
    this.submitting = true;
    this.collectionService.save(this.collection).then(collection => {
      console.log("success");
    }).catch(reason => {
      console.log("fail");
      alert(reason);
    });



  }
}