---
applyTo: '**'
---

# Firefly III API Client - Unimplemented Features

## High Priority Features

### Transaction Management
- [x] View single transaction (GET)
- [x] List all transactions (GET)
- [x] Update existing transaction (PUT)
- [x] Delete transaction (DELETE)
- [x] Search transactions
- [x] Transaction history
- [x] Category integration
- [x] Multi-currency support
- [x] Transaction splits
- [x] Foreign amount handling

### Account Management
- [x] List all accounts (GET)
- [x] View single account (GET)
- [x] Delete account (DELETE)
- [x] Search accounts
- [x] Account roles and types
- [x] Account metadata
- [x] Bank details (IBAN, etc.)
- [x] Account reconciliation

## Medium Priority Features

### Budget System
- [x] Create budget
- [x] Read budget
- [x] Update budget
- [x] Delete budget
- [x] Budget limits
- [x] Budget periods
- [x] Spending tracking
- [x] Auto-budgeting

### Category System
- [x] Create category
- [x] Read category
- [x] Update category
- [x] Delete category
- [x] Category reporting
- [x] Earnings/expenses tracking
- [x] Attachments

### Bill Management
- [x] Create bill
- [x] Read bill
- [x] Update bill
- [x] Delete bill
- [x] Payment tracking
- [x] Rule integration
- [x] Attachments

### Rule System
- [ ] Create rule
- [ ] Read rule
- [ ] Update rule
- [ ] Delete rule
- [ ] Rule groups
- [ ] Testing and triggers
- [ ] Action management

## Lower Priority Features

### Data Management
- [x] Piggy banks
- [x] Tags system
- [x] Charts and reports
- [x] Data import/export
- [x] Bulk operations

### System Features
- [ ] User management
- [ ] Preferences
- [ ] Configuration
- [ ] System information
- [ ] Health checks

### Currency Features
- [x] Currency management
- [x] Exchange rates
- [x] Default currency
- [x] Conversion handling

### Integration Features
- [ ] Webhook management
- [ ] API events
- [ ] External triggers
- [ ] Message handling

## Technical Improvements

### Error Handling
- [x] Detailed error types
- [x] Error categorization
- [x] Rate limit handling
- [x] Retry mechanisms
- [x] Circuit breakers

### Authentication
- [x] Token refresh
- [x] Token validation
- [x] Session management
- [x] Access scopes
- [x] Rate limiting

## Implementation Progress

### Completed Features
- [x] Create account
- [x] Update account balance
- [x] Import single transaction
- [x] Import batch transactions
- [x] Basic error handling
- [x] Bearer token authentication
- [x] View single transaction
- [x] List all transactions
- [x] Update existing transaction
- [x] Delete transaction
- [x] Search transactions
- [x] Category integration
- [x] Multi-currency support
- [x] Transaction splits
- [x] Foreign amount handling
- [x] List all accounts
- [x] View single account
- [x] Delete account
- [x] Search accounts
- [x] Account roles and types
- [x] Account metadata
- [x] Bank details (IBAN, etc.)
- [x] Account reconciliation
- [x] Complete budget system
- [x] Complete category system
- [x] Error handling system
- [x] Currency management
- [x] Complete bill management system
- [x] Tags system
- [x] Charts and reports
- [x] Data import/export
- [x] Bulk operations
- [x] Piggy banks

### In Progress Features
Rule System Implementation

### Next Up
1. Rule System
2. User Management
3. System Features
4. Integration Features
5. Technical Debt Items
6. Documentation Updates

## Importer System

### Core Importers to Implement
1. **Banking System Importers**
   - [ ] Generic CSV Bank Statement Importer
   - [ ] OFX/QFX File Importer
   - [ ] MT940/MT942 Format Importer
   - [ ] CAMT.053 Format Importer

2. **Online Banking APIs**
   - [ ] Plaid Integration
   - [ ] Enable Banking Integration (priority)
   - [ ] TrueLayer Integration
   - [ ] Nordigen/GoCardless Integration
   - [ ] Custom Bank API Adapters

3. **Accounting Software**
   - [ ] QuickBooks Importer
   - [ ] Xero Importer
   - [ ] Wave Accounting Importer
   - [ ] GnuCash File Importer

4. **Investment Platforms**
   - [ ] Interactive Brokers Importer (priority)
   - [ ] TD Ameritrade Importer
   - [ ] Robinhood Importer
   - [ ] Generic Investment CSV Importer

5. **E-commerce Platforms**
   - [ ] PayPal Transaction Importer
   - [ ] Stripe Payment Importer
   - [ ] Amazon Order History Importer
   - [ ] Shopify Transaction Importer

### Features to Implement

1. **Scheduling System**
   - [ ] Cron-based scheduling
   - [ ] Interval-based scheduling
   - [ ] Time zone support
   - [ ] Schedule persistence
   - [ ] Schedule management API
   - [ ] Failure recovery
   - [ ] Schedule history

2. **Data Mapping System**
   - [ ] Template-based mapping
   - [ ] Custom field mapping
   - [ ] Default value handling
   - [ ] Data transformation rules
   - [ ] Validation rules
   - [ ] Error handling
   - [ ] Mapping persistence

3. **Progress Monitoring**
   - [ ] Real-time status updates
   - [ ] Progress persistence
   - [ ] Error tracking
   - [ ] Performance metrics
   - [ ] Resource usage monitoring
   - [ ] Audit logging
   - [ ] Status notifications

4. **Error Recovery**
   - [ ] Automatic retry system
   - [ ] Partial import recovery
   - [ ] Transaction rollback
   - [ ] Error categorization
   - [ ] Recovery strategies
   - [ ] Manual intervention API
   - [ ] Recovery logging

5. **Configuration Management**
   - [ ] Configuration validation
   - [ ] Secure credential storage
   - [ ] Environment variable support
   - [ ] Configuration versioning
   - [ ] Hot reload support
   - [ ] Configuration backup
   - [ ] Migration tools

### Integration Features

1. **Authentication**
   - [ ] OAuth2 support
   - [ ] API key management
   - [ ] Token refresh
   - [ ] Credential encryption
   - [ ] Session management
   - [ ] Rate limit handling
   - [ ] Access logging

2. **Data Validation**
   - [ ] Schema validation
   - [ ] Business rule validation
   - [ ] Currency validation
   - [ ] Date format handling
   - [ ] Amount validation
   - [ ] Duplicate detection
   - [ ] Custom validation rules

3. **Monitoring & Logging**
   - [ ] Structured logging
   - [ ] Performance metrics
   - [ ] Error tracking
   - [ ] Audit trail
   - [ ] Resource monitoring
   - [ ] Health checks
   - [ ] Alert system

4. **Security Features**
   - [ ] Data encryption
   - [ ] Secure communication
   - [ ] Access control
   - [ ] Audit logging
   - [ ] Compliance checks
   - [ ] Data sanitization
   - [ ] Security scanning

### Documentation Needs

1. **Developer Documentation**
   - [ ] Architecture overview
   - [ ] API reference
   - [ ] Integration guides
   - [ ] Security guidelines
   - [ ] Best practices
   - [ ] Troubleshooting guide
   - [ ] Example implementations

2. **User Documentation**
   - [ ] Setup guides
   - [ ] Configuration guides
   - [ ] Usage examples
   - [ ] FAQ
   - [ ] Troubleshooting
   - [ ] Security best practices
   - [ ] Migration guides

3. **Integration Documentation**
   - [ ] API integration guides
   - [ ] Authentication guides
   - [ ] Error handling
   - [ ] Rate limiting
   - [ ] Data formats
   - [ ] Testing guidelines
   - [ ] Production checklist

### Testing Requirements

1. **Unit Tests**
   - [ ] Core functionality
   - [ ] Data validation
   - [ ] Error handling
   - [ ] Configuration
   - [ ] Scheduling
   - [ ] Progress tracking
   - [ ] Resource management

2. **Integration Tests**
   - [ ] API integration
   - [ ] Database operations
   - [ ] Authentication
   - [ ] Error scenarios
   - [ ] Rate limiting
   - [ ] Recovery processes
   - [ ] End-to-end workflows

3. **Performance Tests**
   - [ ] Load testing
   - [ ] Stress testing
   - [ ] Resource usage
   - [ ] Concurrency
   - [ ] Memory leaks
   - [ ] Database performance
   - [ ] API performance 