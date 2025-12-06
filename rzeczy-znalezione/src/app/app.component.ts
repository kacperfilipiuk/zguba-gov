// app.component.ts
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TerritorialUnitsService, TerritorialUnit } from './services/territorial-units.service';

interface FoundItem {
  id: string;
  municipality: {
    name: string;
    type: string;
    contactEmail: string;
  };
  item: {
    name: string;
    category: string;
    date: string;
    location: string;
    status: string;
    description?: string;
  };
  pickup: {
    deadline: number;
    location: string;
    hours?: string;
    contact?: string;
  };
  categories: string[];
  createdAt: string;
}

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  currentStep: number = 1;
  totalSteps: number = 4;
  showErrorModal: boolean = false;
  errorMessages: string[] = [];
  
  form: FormGroup;
  items: FoundItem[] = [];
  
  // Autouzupełnianie
  filteredUnits: TerritorialUnit[] = [];
  showAutocomplete: boolean = false;
  selectedUnitType: TerritorialUnit['type'] | null = null;
  
  categories = [
    { value: 'dokumenty', label: 'Dokumenty' },
    { value: 'telefony', label: 'Telefony i elektronika' },
    { value: 'odziez', label: 'Odzież i akcesoria' },
    { value: 'kluczyki', label: 'Klucze' },
    { value: 'bizesmeny', label: 'Torby i plecaki' },
    { value: 'inne', label: 'Inne' }
  ];

  municipalityTypes = [
    { value: 'powiat', label: 'Powiat' },
    { value: 'gmina', label: 'Gmina' },
    { value: 'miasto', label: 'Miasto' },
    { value: 'wojewodztwo', label: 'Województwo' }
  ];

  statusOptions = [
    { value: 'oczekuje', label: 'Oczekuje na właściciela' },
    { value: 'odebrana', label: 'Odebrana przez właściciela' },
    { value: 'przekazana', label: 'Przekazana organizacji' }
  ];

  constructor(
    private fb: FormBuilder,
    private territorialUnitsService: TerritorialUnitsService
  ) {
    this.form = this.initializeForm();
  }

  ngOnInit(): void {
    this.loadItems();
    this.setupAutocomplete();
  }

  private setupAutocomplete(): void {
    // Nasłuchuj zmian w polu typu samorządu
    this.form.get('municipalityType')?.valueChanges.subscribe(type => {
      this.selectedUnitType = type as TerritorialUnit['type'];
      // Reset nazwy gdy zmienia się typ
      if (this.form.get('municipalityName')?.value) {
        this.onMunicipalityNameInput();
      }
    });
  }

  onMunicipalityNameInput(): void {
    const query = this.form.get('municipalityName')?.value || '';
    
    if (query.length < 2) {
      this.filteredUnits = [];
      this.showAutocomplete = false;
      return;
    }

    this.territorialUnitsService.search(
      query,
      this.selectedUnitType || undefined
    ).then(results => {
      this.filteredUnits = results;
      this.showAutocomplete = results.length > 0;
    });
  }

  selectUnit(unit: TerritorialUnit): void {
    // Generuj sugerowany email
    const suggestedEmail = this.territorialUnitsService.generateContactEmail(unit);
    
    this.form.patchValue({
      municipalityName: unit.name,
      municipalityType: unit.type,
      contactEmail: suggestedEmail
    });
    this.showAutocomplete = false;
    this.filteredUnits = [];
  }

  hideAutocomplete(): void {
    // Opóźnienie aby kliknięcie w sugestię mogło się wykonać
    setTimeout(() => {
      this.showAutocomplete = false;
    }, 200);
  }

  initializeForm(): FormGroup {
    return this.fb.group({
      // Krok 1
      municipalityName: ['', [Validators.required]],
      municipalityType: ['', [Validators.required]],
      contactEmail: ['', [Validators.required, Validators.email]],
      
      // Krok 2
      categories: [[], [Validators.required]],
      
      // Krok 3
      itemName: ['', [Validators.required]],
      itemCategory: ['', [Validators.required]],
      itemDate: ['', [Validators.required]],
      itemLocation: ['', [Validators.required]],
      itemStatus: ['', [Validators.required]],
      itemDescription: [''],
      
      // Krok 4
      storageDeadline: [30, [Validators.required, Validators.min(1), Validators.max(365)]],
      pickupLocation: ['', [Validators.required]],
      pickupHours: [''],
      contactPerson: ['']
    });
  }

  nextStep(): void {
    if (this.validateStep(this.currentStep)) {
      this.currentStep++;
    } else {
      this.showValidationError();
    }
  }

  private showValidationError(): void {
    this.errorMessages = this.getStepErrors(this.currentStep);
    this.showErrorModal = true;
  }

  closeErrorModal(): void {
    this.showErrorModal = false;
    this.errorMessages = [];
  }

  private getStepErrors(step: number): string[] {
    const errors: string[] = [];
    const controls = this.form.controls;

    switch (step) {
      case 1:
        if (!controls['municipalityName'].value?.trim()) {
          errors.push('• Nazwa samorządu jest wymagana');
        }
        if (!controls['municipalityType'].value) {
          errors.push('• Typ samorządu jest wymagany');
        }
        if (!controls['contactEmail'].value?.trim()) {
          errors.push('• Email kontaktowy jest wymagany');
        } else if (!this.isValidEmail(controls['contactEmail'].value)) {
          errors.push('• Email ma nieprawidłowy format (np. kontakt@samorzad.pl)');
        }
        break;
      case 2:
        if (!controls['itemName'].value?.trim()) {
          errors.push('• Nazwa przedmiotu jest wymagana');
        }
        if (!controls['itemCategory'].value) {
          errors.push('• Kategoria przedmiotu jest wymagana');
        }
        if (!controls['itemDate'].value) {
          errors.push('• Data znalezienia jest wymagana');
        }
        if (!controls['itemLocation'].value?.trim()) {
          errors.push('• Miejsce znalezienia jest wymagane');
        }
        if (!controls['itemStatus'].value) {
          errors.push('• Status przedmiotu jest wymagany');
        }
        break;
      case 3:
        if (!controls['storageDeadline'].value || controls['storageDeadline'].value < 1 || controls['storageDeadline'].value > 365) {
          errors.push('• Termin przechowywania musi być między 1 a 365 dni');
        }
        if (!controls['pickupLocation'].value?.trim()) {
          errors.push('• Miejsce odbioru jest wymagane');
        }
        break;
    }

    return errors;
  }

  private isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  prevStep(): void {
    if (this.currentStep > 1) {
      this.currentStep--;
    }
  }

  validateStep(step: number): boolean {
    const controls = this.form.controls;
    
    switch (step) {
      case 1:
        return controls['municipalityName'].valid && 
               controls['municipalityType'].valid && 
               controls['contactEmail'].valid &&
               this.isValidEmail(controls['contactEmail'].value || '');
      case 2:
        return controls['itemName'].valid && 
               controls['itemCategory'].valid && 
               controls['itemDate'].valid && 
               controls['itemLocation'].valid && 
               controls['itemStatus'].valid;
      case 3:
        return controls['storageDeadline'].valid && 
               controls['pickupLocation'].valid;
      default:
        return true;
    }
  }

  submitData(): void {
    const newItem: FoundItem = {
      id: this.generateId(),
      municipality: {
        name: this.form.value.municipalityName,
        type: this.form.value.municipalityType,
        contactEmail: this.form.value.contactEmail
      },
      item: {
        name: this.form.value.itemName,
        category: this.form.value.itemCategory,
        date: this.form.value.itemDate,
        location: this.form.value.itemLocation,
        status: this.form.value.itemStatus,
        description: this.form.value.itemDescription
      },
      pickup: {
        deadline: this.form.value.storageDeadline,
        location: this.form.value.pickupLocation,
        hours: this.form.value.pickupHours,
        contact: this.form.value.contactPerson
      },
      categories: this.form.value.categories,
      createdAt: new Date().toISOString()
    };

    this.items.push(newItem);
    this.saveItems();
    alert('✓ Dane zostały pomyślnie udostępnione!');
    this.resetForm();
    this.currentStep = 1;
  }

  exportJSON(): void {
    const currentItem: FoundItem = {
      id: this.generateId(),
      municipality: {
        name: this.form.value.municipalityName,
        type: this.form.value.municipalityType,
        contactEmail: this.form.value.contactEmail
      },
      item: {
        name: this.form.value.itemName,
        category: this.form.value.itemCategory,
        date: this.form.value.itemDate,
        location: this.form.value.itemLocation,
        status: this.form.value.itemStatus,
        description: this.form.value.itemDescription
      },
      pickup: {
        deadline: this.form.value.storageDeadline,
        location: this.form.value.pickupLocation,
        hours: this.form.value.pickupHours,
        contact: this.form.value.contactPerson
      },
      categories: this.form.value.categories,
      createdAt: new Date().toISOString()
    };

    const dataStr = JSON.stringify([currentItem], null, 2);
    this.downloadFile(dataStr, `rzeczy_znalezione_${Date.now()}.json`, 'application/json');
  }

  exportCSV(): void {
    let csv = 'ID,Samorząd,Typ,Email,Przedmiot,Kategoria,Data,Lokalizacja,Status,Odbiór,Godziny,Osoba Kontaktowa,Data Dodania\n';
    
    const currentItem: FoundItem = {
      id: this.generateId(),
      municipality: {
        name: this.form.value.municipalityName,
        type: this.form.value.municipalityType,
        contactEmail: this.form.value.contactEmail
      },
      item: {
        name: this.form.value.itemName,
        category: this.form.value.itemCategory,
        date: this.form.value.itemDate,
        location: this.form.value.itemLocation,
        status: this.form.value.itemStatus,
        description: this.form.value.itemDescription
      },
      pickup: {
        deadline: this.form.value.storageDeadline,
        location: this.form.value.pickupLocation,
        hours: this.form.value.pickupHours,
        contact: this.form.value.contactPerson
      },
      categories: this.form.value.categories,
      createdAt: new Date().toISOString()
    };

    const row = [
      currentItem.id,
      currentItem.municipality.name,
      currentItem.municipality.type,
      currentItem.municipality.contactEmail,
      currentItem.item.name,
      currentItem.item.category,
      currentItem.item.date,
      currentItem.item.location,
      currentItem.item.status,
      currentItem.pickup.location,
      currentItem.pickup.hours || '',
      currentItem.pickup.contact || '',
      currentItem.createdAt
    ];
    csv += row.map(cell => `"${cell}"`).join(',') + '\n';

    this.downloadFile(csv, `rzeczy_znalezione_${Date.now()}.csv`, 'text/csv');
  }

  private downloadFile(content: string, filename: string, type: string): void {
    const element = document.createElement('a');
    element.setAttribute('href', `data:${type};charset=utf-8,${encodeURIComponent(content)}`);
    element.setAttribute('download', filename);
    element.style.display = 'none';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  }

  private generateId(): string {
    return 'RZ-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
  }

  private saveItems(): void {
    localStorage.setItem('foundItems', JSON.stringify(this.items));
  }

  private loadItems(): void {
    const stored = localStorage.getItem('foundItems');
    if (stored) {
      this.items = JSON.parse(stored);
    }
  }

  private resetForm(): void {
    this.form.reset({
      storageDeadline: 30
    });
  }

  toggleCategory(value: string): void {
    const categoriesControl = this.form.get('categories');
    if (categoriesControl) {
      const currentValue = categoriesControl.value || [];
      const index = currentValue.indexOf(value);
      
      if (index === -1) {
        currentValue.push(value);
      } else {
        currentValue.splice(index, 1);
      }
      
      categoriesControl.setValue([...currentValue]);
    }
  }

  getProgress(): number {
    return (this.currentStep / this.totalSteps) * 100;
  }

  isStepCompleted(step: number): boolean {
    return step < this.currentStep;
  }

  isStepActive(step: number): boolean {
    return step === this.currentStep;
  }

  getFormDataPreview(): string {
    const preview = {
      municipality: {
        name: this.form.value.municipalityName || 'N/A',
        type: this.form.value.municipalityType || 'N/A',
        contactEmail: this.form.value.contactEmail || 'N/A'
      },
      item: {
        name: this.form.value.itemName || 'N/A',
        category: this.form.value.itemCategory || 'N/A',
        date: this.form.value.itemDate || 'N/A',
        location: this.form.value.itemLocation || 'N/A',
        status: this.form.value.itemStatus || 'N/A',
        description: this.form.value.itemDescription || ''
      },
      pickup: {
        deadline: this.form.value.storageDeadline || 30,
        location: this.form.value.pickupLocation || 'N/A',
        hours: this.form.value.pickupHours || '',
        contact: this.form.value.contactPerson || ''
      },
      categories: this.form.value.categories || [],
      createdAt: new Date().toISOString()
    };
    return JSON.stringify(preview, null, 2);
  }
}