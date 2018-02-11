import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CollectionLoaderComponent } from './collection-loader.component';

describe('CollectionLoaderComponent', () => {
  let component: CollectionLoaderComponent;
  let fixture: ComponentFixture<CollectionLoaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CollectionLoaderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CollectionLoaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
